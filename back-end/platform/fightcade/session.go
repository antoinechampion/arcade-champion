package fightcade

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

const retryDelay = 10 * time.Second

type Credentials struct {
	Username string
	Password string
	Cookie   string
}

type Channel struct {
	Name     string `json:"name"`
	GameID   string `json:"gameId"`
	Emulator string `json:"emulator"`
	Players  int    `json:"players"`
	Ranked   bool   `json:"ranked"`
}

type LobbyUser struct {
	Name    string `json:"name"`
	Rank    int    `json:"rank"`
	Playing bool   `json:"playing"`
	Away    bool   `json:"away"`
}

type MatchEvent struct {
	Opponent string
	QuarkID  string
	PlayerID int
	Port     int
	Delay    int
	Ranked   bool
	Token    string
}

var rankNames = []string{"?", "E", "D", "C", "B", "A", "S"}

func RankName(r int) string {
	if r >= 0 && r < len(rankNames) {
		return rankNames[r]
	}
	return fmt.Sprintf("%d", r)
}

func isSuccess(resp map[string]any) bool {
	switch v := resp["result"].(type) {
	case float64:
		return v == 200
	case bool:
		return v
	case string:
		return v == "ok" || v == "success"
	default:
		return false
	}
}

func authenticate(ctx context.Context, client *wsClient, creds Credentials) (map[string]any, error) {
	if creds.Cookie != "" {
		return client.autoLogin(ctx, creds.Cookie)
	}
	return client.login(ctx, creds.Username, creds.Password)
}

func Login(ctx context.Context, creds Credentials) (string, error) {
	client, err := connect(ctx)
	if err != nil {
		return "", err
	}
	defer client.close()

	resp, err := authenticate(ctx, client, creds)
	if err != nil {
		return "", err
	}
	if !isSuccess(resp) {
		return "", fmt.Errorf("login failed: %v", resp["error"])
	}

	cookie, _ := resp["cookie"].(string)
	if cookie == "" {
		return "", fmt.Errorf("login succeeded but server returned no cookie")
	}
	return cookie, nil
}

func Search(ctx context.Context, creds Credentials, query string) ([]Channel, error) {
	client, err := connect(ctx)
	if err != nil {
		return nil, err
	}
	defer client.close()

	resp, err := authenticate(ctx, client, creds)
	if err != nil {
		return nil, err
	}
	if !isSuccess(resp) {
		return nil, fmt.Errorf("login failed: %v", resp["error"])
	}

	result, err := client.searchChannels(ctx, query, 0)
	if err != nil {
		return nil, err
	}

	raw, _ := result["channels"].([]any)
	channels := make([]Channel, 0, len(raw))
	for _, r := range raw {
		ch, _ := r.(map[string]any)
		name, _ := ch["name"].(string)
		gameid, _ := ch["gameid"].(string)
		emu, _ := ch["emulator"].(string)
		players := 0
		if p, ok := ch["clients"].(float64); ok {
			players = int(p)
		}
		ranked, _ := ch["ranked"].(bool)
		channels = append(channels, Channel{
			Name:     name,
			GameID:   gameid,
			Emulator: emu,
			Players:  players,
			Ranked:   ranked,
		})
	}
	return channels, nil
}

func resolveChannel(ctx context.Context, client *wsClient, game string) (string, error) {
	result, err := client.searchChannels(ctx, game, 0)
	if err != nil {
		return "", err
	}
	channels, _ := result["channels"].([]any)
	for _, raw := range channels {
		ch, _ := raw.(map[string]any)
		if gid, _ := ch["gameid"].(string); gid == game {
			name, _ := ch["name"].(string)
			return name, nil
		}
	}
	for _, raw := range channels {
		ch, _ := raw.(map[string]any)
		if name, _ := ch["name"].(string); name == game {
			return name, nil
		}
	}
	if len(channels) > 0 {
		ch, _ := channels[0].(map[string]any)
		name, _ := ch["name"].(string)
		return name, nil
	}
	return "", fmt.Errorf("no channel found for %q", game)
}

func parseLobbyUsers(users []any, localUser string) ([]LobbyUser, int) {
	var parsed []LobbyUser
	myRank := 0
	for _, raw := range users {
		u, _ := raw.(map[string]any)
		name, _ := u["name"].(string)
		rank := 0
		if r, ok := u["rank"].(float64); ok {
			rank = int(r)
		}
		playing := u["playing"] != nil
		away, _ := u["away"].(bool)
		if ca, _ := u["channel_away"].(bool); ca {
			away = true
		}
		if name == localUser {
			myRank = rank
		}
		parsed = append(parsed, LobbyUser{Name: name, Rank: rank, Playing: playing, Away: away})
	}
	return parsed, myRank
}

func sortByRankDistance(users []LobbyUser, myRank int, myName string) []LobbyUser {
	var candidates []LobbyUser
	for _, u := range users {
		if u.Name == myName || u.Playing || u.Away {
			continue
		}
		candidates = append(candidates, u)
	}
	for i := 1; i < len(candidates); i++ {
		key := candidates[i]
		keyDist := math.Abs(float64(key.Rank - myRank))
		j := i - 1
		for j >= 0 && math.Abs(float64(candidates[j].Rank-myRank)) > keyDist {
			candidates[j+1] = candidates[j]
			j--
		}
		candidates[j+1] = key
	}
	return candidates
}

func Lobby(ctx context.Context, creds Credentials, game string) (*MatchEvent, error) {
	client, err := connect(ctx)
	if err != nil {
		return nil, err
	}
	defer client.close()

	resp, err := authenticate(ctx, client, creds)
	if err != nil {
		return nil, err
	}
	if !isSuccess(resp) {
		return nil, fmt.Errorf("login failed: %v", resp["error"])
	}

	user, _ := resp["user"].(map[string]any)
	token, _ := user["token"].(string)
	username, _ := user["name"].(string)

	channelname, err := resolveChannel(ctx, client, game)
	if err != nil {
		return nil, err
	}

	joinResp, err := client.joinChannel(ctx, channelname)
	if err != nil {
		return nil, err
	}
	if !isSuccess(joinResp) {
		return nil, fmt.Errorf("failed to join %s: %v", channelname, joinResp["error"])
	}

	emulator, _ := joinResp["emulator"].(string)
	if emulator == "" {
		emulator = "fbneo"
	}
	gameid, _ := joinResp["gameid"].(string)
	if gameid == "" {
		gameid = game
	}
	isRanked, _ := joinResp["ranked"].(bool)

	usersRaw, _ := joinResp["users"].([]any)
	lobbyUsers, myRank := parseLobbyUsers(usersRaw, username)

	matchCh := make(chan *MatchEvent, 1)
	challenged := make(chan struct{})
	var closeOnce sync.Once

	client.on("start", func(msg map[string]any) {
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
		matchToken, _ := msg["token"].(string)
		if matchToken == "" {
			matchToken = token
		}

		url := buildMatchURL(emulator, gameid, quarkid, playerid, port, delay, ranked, matchToken)
		_ = openURL(url)

		matchCh <- &MatchEvent{
			Opponent: opponent,
			QuarkID:  quarkid,
			PlayerID: playerid,
			Port:     port,
			Delay:    delay,
			Ranked:   ranked,
			Token:    matchToken,
		}
	})

	client.on("challenge", func(msg map[string]any) {
		select {
		case <-challenged:
			return
		default:
		}
		userObj, _ := msg["user"].(map[string]any)
		opponent, _ := userObj["name"].(string)
		channel, _ := msg["channelname"].(string)
		cidFloat, _ := msg["challengeid"].(float64)
		cid := int(cidFloat)
		ranked, _ := msg["ranked"].(bool)
		go func() {
			_, _ = client.acceptChallenge(ctx, opponent, channel, cid, ranked)
		}()
	})

	go func() {
		candidates := sortByRankDistance(lobbyUsers, myRank, username)
		challengeID := 1
		tried := map[string]bool{}

		acceptCh := make(chan bool, 1)
		client.on("accept", func(msg map[string]any) {
			closeOnce.Do(func() { close(challenged) })
			select {
			case acceptCh <- true:
			default:
			}
		})
		client.on("reject", func(msg map[string]any) {
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
			resp, err := client.challengeUser(ctx, target.Name, channelname, challengeID, isRanked)
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
	}()

	defer func() {
		_ = client.leaveChannel(channelname)
	}()

	select {
	case match := <-matchCh:
		return match, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-client.done:
		return nil, fmt.Errorf("connection closed")
	}
}

func Play(emulator, gameid string) error {
	return openURL(buildPlayURL(emulator, gameid))
}

func Training(emulator, gameid string) error {
	return openURL(buildTrainingURL(emulator, gameid))
}
