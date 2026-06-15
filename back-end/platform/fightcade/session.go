package fightcade

import (
	"context"
	"fmt"
	"log"
	"time"
)

// retryDelay paces how long we wait before challenging the next candidate while
// earlier challenges remain outstanding. A var so tests can shrink it.
var retryDelay = 15 * time.Second

// acceptDelay is how long we wait before accepting an incoming challenge, giving
// the challenge loop time to win via an outgoing accept instead. A var so tests can shrink it.
var acceptDelay = 3 * time.Second

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
	Vping   int    `json:"vping"` // connection-bar level shown in the lobby UI (higher = better)
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
		log.Printf("[fightcade] authenticating with cookie")
		return client.autoLogin(ctx, creds.Cookie)
	}
	log.Printf("[fightcade] authenticating with username=%s", creds.Username)
	return client.login(ctx, creds.Username, creds.Password)
}

func Login(ctx context.Context, creds Credentials) (string, error) {
	log.Printf("[fightcade] Login: connecting to server")
	client, err := connect(ctx)
	if err != nil {
		log.Printf("[fightcade] Login: connection failed: %v", err)
		return "", err
	}
	defer client.close()

	resp, err := authenticate(ctx, client, creds)
	if err != nil {
		log.Printf("[fightcade] Login: auth error: %v", err)
		return "", err
	}
	if !isSuccess(resp) {
		log.Printf("[fightcade] Login: auth rejected: %v", resp["error"])
		return "", fmt.Errorf("login failed: %v", resp["error"])
	}

	cookie, _ := resp["cookie"].(string)
	if cookie == "" {
		log.Printf("[fightcade] Login: no cookie in response")
		return "", fmt.Errorf("login succeeded but server returned no cookie")
	}
	log.Printf("[fightcade] Login: success")
	return cookie, nil
}

type SearchResult struct {
	Channels []Channel
	Cookie   string
}

func Search(ctx context.Context, creds Credentials, query string) (SearchResult, error) {
	log.Printf("[fightcade] Search: query=%q", query)
	client, err := connect(ctx)
	if err != nil {
		log.Printf("[fightcade] Search: connection failed: %v", err)
		return SearchResult{}, err
	}
	defer client.close()

	resp, err := authenticate(ctx, client, creds)
	if err != nil {
		log.Printf("[fightcade] Search: auth error: %v", err)
		return SearchResult{}, err
	}
	if !isSuccess(resp) {
		log.Printf("[fightcade] Search: auth rejected: %v", resp["error"])
		return SearchResult{}, fmt.Errorf("login failed: %v", resp["error"])
	}

	cookie, _ := resp["cookie"].(string)

	result, err := client.searchChannels(ctx, query, 0)
	if err != nil {
		log.Printf("[fightcade] Search: searchChannels failed: %v", err)
		return SearchResult{}, err
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
	log.Printf("[fightcade] Search: found %d channels", len(channels))
	return SearchResult{Channels: channels, Cookie: cookie}, nil
}

func Lobby(ctx context.Context, creds Credentials, game string) (*MatchEvent, error) {
	log.Printf("[fightcade] Lobby: starting for game=%q", game)
	client, err := connect(ctx)
	if err != nil {
		log.Printf("[fightcade] Lobby: connection failed: %v", err)
		return nil, err
	}
	defer client.close()

	resp, err := authenticate(ctx, client, creds)
	if err != nil {
		log.Printf("[fightcade] Lobby: auth error: %v", err)
		return nil, err
	}
	if !isSuccess(resp) {
		log.Printf("[fightcade] Lobby: auth rejected: %v", resp["error"])
		return nil, fmt.Errorf("login failed: %v", resp["error"])
	}

	user, _ := resp["user"].(map[string]any)
	token, _ := user["token"].(string)
	username, _ := user["name"].(string)
	log.Printf("[fightcade] Lobby: logged in as %q", username)

	reportStatus(token)

	channelname, err := resolveChannel(ctx, client, game)
	if err != nil {
		log.Printf("[fightcade] Lobby: resolveChannel failed: %v", err)
		return nil, err
	}
	log.Printf("[fightcade] Lobby: resolved channel=%q", channelname)

	joinResp, err := client.joinChannel(ctx, channelname)
	if err != nil {
		log.Printf("[fightcade] Lobby: joinChannel error: %v", err)
		return nil, err
	}
	if !isSuccess(joinResp) {
		log.Printf("[fightcade] Lobby: joinChannel rejected: %v", joinResp["error"])
		return nil, fmt.Errorf("failed to join %s: %v", channelname, joinResp["error"])
	}
	defer func() { _ = client.leaveChannel(channelname) }()

	if err := client.setNotAway(channelname); err != nil {
		log.Printf("[fightcade] Lobby: setNotAway error: %v", err)
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
	log.Printf("[fightcade] Lobby: channel=%q emulator=%s gameid=%s ranked=%v", channelname, emulator, gameid, isRanked)

	usersRaw, _ := joinResp["users"].([]any)
	lobbyUsers, myRank := parseLobbyUsers(usersRaw, username)
	log.Printf("[fightcade] Lobby: %d users in lobby, my rank=%s(%d)", len(lobbyUsers), RankName(myRank), myRank)
	for _, u := range lobbyUsers {
		if u.Name == username || u.Playing || u.Away {
			continue
		}
		log.Printf("[fightcade] Lobby:   user=%s rank=%s(%d) connection=%d bars", u.Name, RankName(u.Rank), u.Rank, u.Vping)
	}

	mm := &matchmaker{client: client, config: lobbyConfig{
		channelName: channelname,
		emulator:    emulator,
		gameID:      gameid,
		ranked:      isRanked,
		username:    username,
		token:       token,
		users:       lobbyUsers,
		myRank:      myRank,
	}}
	return mm.run(ctx)
}

func resolveChannel(ctx context.Context, client *wsClient, game string) (string, error) {
	log.Printf("[fightcade] resolveChannel: searching for game=%q", game)
	result, err := client.searchChannels(ctx, game, 0)
	if err != nil {
		return "", err
	}
	channels, _ := result["channels"].([]any)
	log.Printf("[fightcade] resolveChannel: got %d channels", len(channels))
	for _, raw := range channels {
		ch, _ := raw.(map[string]any)
		gid, _ := ch["gameid"].(string)
		name, _ := ch["name"].(string)
		log.Printf("[fightcade] resolveChannel:   name=%q gameid=%q", name, gid)
		if gid == game {
			log.Printf("[fightcade] resolveChannel: matched by gameid -> %q", name)
			return name, nil
		}
	}
	for _, raw := range channels {
		ch, _ := raw.(map[string]any)
		if name, _ := ch["name"].(string); name == game {
			log.Printf("[fightcade] resolveChannel: matched by name -> %q", name)
			return name, nil
		}
	}
	if len(channels) > 0 {
		ch, _ := channels[0].(map[string]any)
		name, _ := ch["name"].(string)
		log.Printf("[fightcade] resolveChannel: no exact match, using first result -> %q", name)
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
		vping := 0
		if v, ok := u["vping"].(float64); ok {
			vping = int(v)
		}
		playing := false
		if p, ok := u["playing"].(bool); ok {
			playing = p
		} else if _, ok := u["playing"].(map[string]any); ok {
			playing = true
		}
		away, _ := u["away"].(bool)
		if ca, _ := u["channel_away"].(bool); ca {
			away = true
		}
		if name == localUser {
			myRank = rank
		}
		parsed = append(parsed, LobbyUser{Name: name, Rank: rank, Vping: vping, Playing: playing, Away: away})
	}
	return parsed, myRank
}

func Play(emulator, gameid string) error {
	return openURL(buildPlayURL(emulator, gameid))
}

func Training(emulator, gameid string) error {
	return openURL(buildTrainingURL(emulator, gameid))
}
