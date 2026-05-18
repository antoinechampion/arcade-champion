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
	client *wsClient
	config lobbyConfig
}


func (m *matchmaker) run(ctx context.Context) (*MatchEvent, error) {
	log.Printf("[fightcade] matchmaker: starting (channel=%s ranked=%v debugMode=%v)", m.config.channelName, m.config.ranked, debugMode)
	if debugMode {
		log.Printf("[fightcade] matchmaker: DEBUG MODE — will only challenge target=%q", debugTargetPlayer)
	}

	matchCh := make(chan *MatchEvent, 1)
	challenged := make(chan struct{})
	var closeOnce sync.Once

	m.client.on("start", func(msg map[string]any) {
		event := parseStartEvent(msg, m.config.token)
		log.Printf("[fightcade] matchmaker: START event — opponent=%s quarkid=%s playerid=%d port=%d delay=%d ranked=%v",
			event.Opponent, event.QuarkID, event.PlayerID, event.Port, event.Delay, event.Ranked)
		launchGame(m.config.emulator, m.config.gameID, event)
		matchCh <- event
	})

	m.client.on("challenge", func(msg map[string]any) {
		select {
		case <-challenged:
			log.Printf("[fightcade] matchmaker: incoming challenge ignored (already matched)")
			return
		default:
		}
		userObj, _ := msg["user"].(map[string]any)
		challenger, _ := userObj["name"].(string)
		log.Printf("[fightcade] matchmaker: incoming challenge from %q", challenger)
		m.acceptIncoming(ctx, msg)
	})

	go m.challengeLoop(ctx, challenged, &closeOnce)

	select {
	case match := <-matchCh:
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

func (m *matchmaker) challengeLoop(ctx context.Context, challenged chan struct{}, closeOnce *sync.Once) {
	candidates := sortByRankDistance(m.config.users, m.config.myRank, m.config.username)
	log.Printf("[fightcade] challengeLoop: %d candidates sorted by rank distance (my rank=%s(%d)):",
		len(candidates), RankName(m.config.myRank), m.config.myRank)
	for i, c := range candidates {
		dist := math.Abs(float64(c.Rank - m.config.myRank))
		log.Printf("[fightcade] challengeLoop:   #%d %s rank=%s(%d) distance=%.0f", i+1, c.Name, RankName(c.Rank), c.Rank, dist)
	}

	if debugMode {
		var filtered []LobbyUser
		for _, c := range candidates {
			if c.Name == debugTargetPlayer {
				filtered = append(filtered, c)
			}
		}
		if len(filtered) == 0 {
			log.Printf("[fightcade] challengeLoop: DEBUG MODE — target %q not found in lobby, waiting for incoming challenges only", debugTargetPlayer)
		} else {
			log.Printf("[fightcade] challengeLoop: DEBUG MODE — narrowed candidates to target %q only", debugTargetPlayer)
		}
		candidates = filtered
	}

	challengeID := 1
	tried := map[string]bool{}

	acceptCh := make(chan bool, 1)
	m.client.on("accept", func(msg map[string]any) {
		userObj, _ := msg["user"].(map[string]any)
		accepter, _ := userObj["name"].(string)
		log.Printf("[fightcade] challengeLoop: challenge ACCEPTED by %q", accepter)
		closeOnce.Do(func() { close(challenged) })
		select {
		case acceptCh <- true:
		default:
		}
	})
	m.client.on("reject", func(msg map[string]any) {
		userObj, _ := msg["user"].(map[string]any)
		rejecter, _ := userObj["name"].(string)
		log.Printf("[fightcade] challengeLoop: challenge REJECTED by %q", rejecter)
		select {
		case acceptCh <- false:
		default:
		}
	})

	for {
		var target *LobbyUser
		for i := range candidates {
			if !tried[candidates[i].Name] {
				target = &candidates[i]
				break
			}
		}
		if target == nil {
			log.Printf("[fightcade] challengeLoop: no more candidates to challenge")
			return
		}

		tried[target.Name] = true
		log.Printf("[fightcade] challengeLoop: challenging %q (rank=%s(%d) challengeID=%d)",
			target.Name, RankName(target.Rank), target.Rank, challengeID)
		resp, err := m.client.challengeUser(ctx, target.Name, m.config.channelName, challengeID, m.config.ranked)
		challengeID++
		if err != nil {
			log.Printf("[fightcade] challengeLoop: challengeUser error: %v", err)
			continue
		}
		if !isSuccess(resp) {
			log.Printf("[fightcade] challengeLoop: challengeUser failed: %v", resp)
			continue
		}
		log.Printf("[fightcade] challengeLoop: challenge sent to %q, waiting for response (timeout=%s)", target.Name, retryDelay)
		select {
		case accepted := <-acceptCh:
			if accepted {
				log.Printf("[fightcade] challengeLoop: match confirmed with %q", target.Name)
				return
			}
			log.Printf("[fightcade] challengeLoop: %q rejected, moving to next candidate", target.Name)
		case <-time.After(retryDelay):
			log.Printf("[fightcade] challengeLoop: timeout waiting for %q, moving to next candidate", target.Name)
		case <-ctx.Done():
			log.Printf("[fightcade] challengeLoop: context cancelled")
			return
		}
	}
}

func (m *matchmaker) acceptIncoming(ctx context.Context, msg map[string]any) {
	userObj, _ := msg["user"].(map[string]any)
	opponent, _ := userObj["name"].(string)
	channel, _ := msg["channelname"].(string)
	cidFloat, _ := msg["challengeid"].(float64)
	cid := int(cidFloat)
	ranked, _ := msg["ranked"].(bool)

	if debugMode && opponent != debugTargetPlayer {
		log.Printf("[fightcade] acceptIncoming: DEBUG MODE — ignoring challenge from %q (only accepting from %q)", opponent, debugTargetPlayer)
		return
	}

	log.Printf("[fightcade] acceptIncoming: accepting challenge from %q (channel=%s challengeID=%d ranked=%v)", opponent, channel, cid, ranked)
	go func() {
		_, err := m.client.acceptChallenge(ctx, opponent, channel, cid, ranked)
		if err != nil {
			log.Printf("[fightcade] acceptIncoming: acceptChallenge error: %v", err)
		}
	}()
}

func parseStartEvent(msg map[string]any, fallbackToken string) *MatchEvent {
	userObj, _ := msg["user"].(map[string]any)
	opponent, _ := userObj["name"].(string)
	quarkid, _ := msg["quarkid"].(string)
	playerid := int(msg["playerid"].(float64))
	port := int(msg["port"].(float64))
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

func launchGame(emulator, gameID string, event *MatchEvent) {
	url := buildMatchURL(emulator, gameID, event.QuarkID, event.PlayerID, event.Port, event.Delay, event.Ranked, event.Token)
	log.Printf("[fightcade] launchGame: opening url=%s", url)
	if err := openURL(url); err != nil {
		log.Printf("[fightcade] launchGame: openURL error: %v", err)
	}
}

func sortByRankDistance(users []LobbyUser, myRank int, myName string) []LobbyUser {
	var candidates []LobbyUser
	for _, u := range users {
		if u.Name == myName || u.Playing || u.Away {
			continue
		}
		candidates = append(candidates, u)
	}
	slices.SortFunc(candidates, func(a, b LobbyUser) int {
		da := math.Abs(float64(a.Rank - myRank))
		db := math.Abs(float64(b.Rank - myRank))
		return cmp.Compare(da, db)
	})
	return candidates
}
