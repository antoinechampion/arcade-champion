package fightcade

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"math"
	"slices"
	"sync"
	"time"
)

type lobbyConfig struct {
	channelName string
	emulator    string
	gameID      string
	ranked      bool
	username    string
	token       string
	users       []LobbyUser
	myRank      int
}

type matchmaker struct {
	client  *wsClient
	config  lobbyConfig
	matched string // opponent we matched with; set once, read when cancelling pending challenges
}

// run drives matchmaking until a game starts, the context is cancelled, or the
// connection drops. A match can be reached two ways — an opponent accepts one of
// our outgoing challenges, or we accept an incoming one — but both converge on the
// server's authoritative "start" event, which is the single signal that wins.
func (m *matchmaker) run(ctx context.Context) (*MatchEvent, error) {
	log.Printf("[fightcade] matchmaker: starting (channel=%s ranked=%v debugMode=%v)", m.config.channelName, m.config.ranked, debugMode)
	if debugMode {
		log.Printf("[fightcade] matchmaker: DEBUG MODE — dry run, no challenge or accept requests will be sent")
	}

	// matchCtx is cancelled the moment a match is won; the challenge loop watches
	// it to stop challenging and cancel any still-pending challenges.
	matchCtx, cancelMatch := context.WithCancel(ctx)
	defer cancelMatch()

	matchmakingCh := make(chan *MatchEvent, 1)
	rejectCh := make(chan string, len(m.config.users)+1)
	var matchOnce sync.Once

	m.client.on("start", func(msg map[string]any) {
		event := parseStartEvent(msg, m.config.token)
		log.Printf("[fightcade] matchmaker: START event — opponent=%s quarkid=%s playerid=%d port=%d delay=%d ranked=%v",
			event.Opponent, event.QuarkID, event.PlayerID, event.Port, event.Delay, event.Ranked)
		// sync.Once makes it safe against duplicate start events (which would otherwise deadlock recvLoop on
		// the second matchCh send) and double game launches.
		matchOnce.Do(func() {
			m.matched = event.Opponent
			launchGame(m.config.emulator, m.config.gameID, event)
			matchmakingCh <- event
			cancelMatch()
		})
	})

	m.client.on("challenge", func(msg map[string]any) {
		if matchCtx.Err() != nil {
			log.Printf("[fightcade] matchmaker: incoming challenge ignored (already matched)")
			return
		}
		log.Printf("[fightcade] matchmaker: incoming challenge from %q", userName(msg))
		m.acceptIncoming(matchCtx, msg)
	})

	m.client.on("accept", func(msg map[string]any) {
		log.Printf("[fightcade] matchmaker: challenge ACCEPTED by %q", userName(msg))
	})

	m.client.on("reject", func(msg map[string]any) {
		rejecter := userName(msg)
		log.Printf("[fightcade] matchmaker: challenge REJECTED by %q", rejecter)
		select {
		case rejectCh <- rejecter:
		default:
		}
	})

	go m.challengeLoop(matchCtx, rejectCh)

	select {
	case match := <-matchmakingCh:
		log.Printf("[fightcade] matchmaker: match found — opponent=%s", match.Opponent)
		return match, nil
	case <-ctx.Done():
		log.Printf("[fightcade] matchmaker: context cancelled")
		return nil, ctx.Err()
	case <-m.client.done:
		log.Printf("[fightcade] matchmaker: connection closed unexpectedly")
		return nil, fmt.Errorf("connection closed")
	}
}

// challengeNext sends a challenge to the next untried candidate. Returns false
// when no candidates remain.
func (m *matchmaker) challengeNext(ctx context.Context, candidates []LobbyUser, next *int, pending map[string]int) bool {
	if *next >= len(candidates) {
		return false
	}
	target := candidates[*next]
	challengeID := *next + 1
	*next++
	log.Printf("[fightcade] challengeLoop: challenging %q (rank=%s(%d) challengeID=%d)",
		target.Name, RankName(target.Rank), target.Rank, challengeID)
	resp, err := m.client.challengeUser(ctx, target.Name, m.config.channelName, challengeID, m.config.ranked)
	if err != nil {
		log.Printf("[fightcade] challengeLoop: challengeUser error: %v", err)
		return true
	}
	if !isSuccess(resp) {
		log.Printf("[fightcade] challengeLoop: challengeUser failed: %v", resp)
		return true
	}
	pending[target.Name] = challengeID
	log.Printf("[fightcade] challengeLoop: challenge sent to %q, %d pending", target.Name, len(pending))
	return true
}

// challengeLoop challenges candidates one at a time, pacing a new challenge every
// retryDelay (or sooner when one is rejected). Challenges stay live, so several can
// be outstanding at once; the first opponent to accept wins. When ctx is cancelled
// — a match was won — every still-pending challenge except the winner's is cancelled.
func (m *matchmaker) challengeLoop(ctx context.Context, rejectCh <-chan string) {
	candidates := sortCandidates(m.config.users, m.config.myRank, m.config.username)

	pending := map[string]int{} // opponent name -> challengeid
	next := 0

	m.challengeNext(ctx, candidates, &next, pending)

	for {
		if next >= len(candidates) && len(pending) == 0 {
			log.Printf("[fightcade] challengeLoop: no candidates left and nothing pending, only incoming challenges can match now")
			return
		}
		select {
		case <-ctx.Done():
			m.cancelPending(pending)
			return
		case name := <-rejectCh:
			delete(pending, name)
			m.challengeNext(ctx, candidates, &next, pending)
		case <-time.After(retryDelay):
			m.challengeNext(ctx, candidates, &next, pending)
		}
	}
}

// cancelPending cancels every pending challenge except the one we matched with.
// It runs after ctx is cancelled, so it uses a fresh context for the requests.
func (m *matchmaker) cancelPending(pending map[string]int) {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()
	for name, cid := range pending {
		if name == m.matched {
			continue
		}
		log.Printf("[fightcade] challengeLoop: cancelling pending challenge to %q (challengeID=%d)", name, cid)
		if _, err := m.client.cancelChallenge(ctx, name, m.config.channelName, cid); err != nil {
			log.Printf("[fightcade] challengeLoop: cancelChallenge error: %v", err)
		}
	}
}

func (m *matchmaker) acceptIncoming(ctx context.Context, msg map[string]any) {
	opponent := userName(msg)
	channel, _ := msg["channelname"].(string)
	cidFloat, _ := msg["challengeid"].(float64)
	cid := int(cidFloat)
	ranked, _ := msg["ranked"].(bool)

	log.Printf("[fightcade] acceptIncoming: accepting challenge from %q (channel=%s challengeID=%d ranked=%v) in 3s", opponent, channel, cid, ranked)
	go func() {
		select {
		case <-ctx.Done():
			log.Printf("[fightcade] acceptIncoming: context cancelled before delay elapsed, not accepting %q", opponent)
			return
		case <-time.After(acceptDelay):
		}
		if ctx.Err() != nil {
			log.Printf("[fightcade] acceptIncoming: context cancelled after delay, not accepting %q", opponent)
			return
		}
		_, err := m.client.acceptChallenge(ctx, opponent, channel, cid, ranked)
		if err != nil {
			log.Printf("[fightcade] acceptIncoming: acceptChallenge error: %v", err)
		}
	}()
}

func userName(msg map[string]any) string {
	user, _ := msg["user"].(map[string]any)
	name, _ := user["name"].(string)
	return name
}

func parseStartEvent(msg map[string]any, fallbackToken string) *MatchEvent {
	opponent := userName(msg)
	quarkid, _ := msg["quarkid"].(string)
	playeridF, _ := msg["playerid"].(float64)
	portF, _ := msg["port"].(float64)
	playerid := int(playeridF)
	port := int(portF)
	delay := 0
	if d, ok := msg["delay"].(float64); ok {
		delay = int(d)
	}
	ranked, _ := msg["ranked"].(bool)
	token, _ := msg["token"].(string)
	if token == "" {
		token = fallbackToken
	}
	return &MatchEvent{
		Opponent: opponent,
		QuarkID:  quarkid,
		PlayerID: playerid,
		Port:     port,
		Delay:    delay,
		Ranked:   ranked,
		Token:    token,
	}
}

// launchGame is a var so tests can stub the side-effecting URL open.
var launchGame = func(emulator, gameID string, event *MatchEvent) {
	url := buildMatchURL(emulator, gameID, event.QuarkID, event.PlayerID, event.Port, event.Delay, event.Ranked, event.Token)
	log.Printf("[fightcade] launchGame: opening url=%s", url)
	if err := openURL(url); err != nil {
		log.Printf("[fightcade] launchGame: openURL error: %v", err)
	}
}

// vping is the lobby connection-bar level (higher = better)
const (
	minConnection  = 1 // drop anyone with fewer bars
	goodConnection = 3 // challenged before any weaker connection
)

// sortCandidates drops poorly-connected players, then orders the rest so a good
// connection always beats a closer rank
func sortCandidates(users []LobbyUser, myRank int, myName string) []LobbyUser {
	var candidates []LobbyUser
	for _, u := range users {
		if u.Name == myName || u.Playing || u.Away || u.Vping <= minConnection {
			continue
		}
		candidates = append(candidates, u)
	}
	slices.SortFunc(candidates, func(a, b LobbyUser) int {
		ga := a.Vping >= goodConnection
		gb := b.Vping >= goodConnection
		if ga != gb {
			if ga {
				return -1
			}
			return 1
		}
		da := math.Abs(float64(a.Rank - myRank))
		db := math.Abs(float64(b.Rank - myRank))
		return cmp.Compare(da, db)
	})

	log.Printf("[fightcade] challengeLoop: %d candidates found (my rank=%s(%d)):",
		len(candidates), RankName(myRank), myRank)
	for i, c := range candidates {
		dist := math.Abs(float64(c.Rank - myRank))
		log.Printf("- #%d %s rank=%s(%d) vping=%d distance=%.0f", i+1, c.Name, RankName(c.Rank), c.Rank, c.Vping, dist)
	}

	return candidates
}
