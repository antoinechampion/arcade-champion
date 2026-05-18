package fightcade

import (
	"cmp"
	"context"
	"fmt"
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
	matchCh := make(chan *MatchEvent, 1)
	challenged := make(chan struct{})
	var closeOnce sync.Once

	m.client.on("start", func(msg map[string]any) {
		event := parseStartEvent(msg, m.config.token)
		launchGame(m.config.emulator, m.config.gameID, event)
		matchCh <- event
	})

	m.client.on("challenge", func(msg map[string]any) {
		select {
		case <-challenged:
			return
		default:
		}
		m.acceptIncoming(ctx, msg)
	})

	go m.challengeLoop(ctx, challenged, &closeOnce)

	select {
	case match := <-matchCh:
		return match, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-m.client.done:
		return nil, fmt.Errorf("connection closed")
	}
}

func (m *matchmaker) challengeLoop(ctx context.Context, challenged chan struct{}, closeOnce *sync.Once) {
	candidates := sortByRankDistance(m.config.users, m.config.myRank, m.config.username)
	challengeID := 1
	tried := map[string]bool{}

	acceptCh := make(chan bool, 1)
	m.client.on("accept", func(msg map[string]any) {
		closeOnce.Do(func() { close(challenged) })
		select {
		case acceptCh <- true:
		default:
		}
	})
	m.client.on("reject", func(msg map[string]any) {
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
			return
		}

		tried[target.Name] = true
		resp, err := m.client.challengeUser(ctx, target.Name, m.config.channelName, challengeID, m.config.ranked)
		challengeID++
		if err == nil && isSuccess(resp) {
			select {
			case accepted := <-acceptCh:
				if accepted {
					return
				}
			case <-time.After(retryDelay):
			case <-ctx.Done():
				return
			}
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
	go func() {
		_, _ = m.client.acceptChallenge(ctx, opponent, channel, cid, ranked)
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
	_ = openURL(url)
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
