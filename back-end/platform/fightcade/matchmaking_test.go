package fightcade

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// TestMatchmaker_FirstAcceptWinsAndCancelsRest challenges several candidates, has
// the server start a match with one of them, and verifies the matchmaker returns
// that opponent and cancels the other still-pending challenge.
func TestMatchmaker_FirstAcceptWinsAndCancelsRest(t *testing.T) {
	prevDelay, prevLaunch := retryDelay, launchGame
	retryDelay = 5 * time.Millisecond
	launchGame = func(emulator, gameID string, event *MatchEvent) {}
	t.Cleanup(func() { retryDelay, launchGame = prevDelay, prevLaunch })

	var mu sync.Mutex
	challenged := map[string]int{} // name -> challengeid
	cancelled := []string{}
	started := false

	client := startTestServer(t, func(conn *websocket.Conn, msg map[string]any) {
		req, _ := msg["req"].(string)
		name, _ := msg["username"].(string)
		idx := msg["requestIdx"]

		mu.Lock()
		defer mu.Unlock()
		switch req {
		case "challenge":
			cid, _ := msg["challengeid"].(float64)
			challenged[name] = int(cid)
			reply(conn, map[string]any{"req": req, "result": 200, "requestIdx": idx})
			// Once both candidates have live challenges, "winner" accepts: the
			// server pushes the authoritative start event.
			if len(challenged) == 2 && !started {
				started = true
				reply(conn, map[string]any{
					"req": "start", "user": map[string]any{"name": "winner"},
					"quarkid": "1-2", "playerid": float64(0), "port": float64(6000),
					"ranked": false, "token": "tok",
				})
			}
		case "cancel":
			cancelled = append(cancelled, name)
			reply(conn, map[string]any{"req": req, "result": 200, "requestIdx": idx})
		default:
			reply(conn, map[string]any{"req": req, "result": 200, "requestIdx": idx})
		}
	})

	mm := &matchmaker{client: client, config: lobbyConfig{
		channelName: "chan",
		username:    "me",
		myRank:      3,
		users: []LobbyUser{
			{Name: "winner", Rank: 3},
			{Name: "loser", Rank: 4},
		},
	}}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	match, err := mm.run(ctx)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if match.Opponent != "winner" {
		t.Errorf("matched opponent = %q, want winner", match.Opponent)
	}

	// cancelPending runs in the challenge loop after the match is won; give it a moment.
	deadline := time.Now().Add(time.Second)
	for {
		mu.Lock()
		done := len(cancelled) == 1 && cancelled[0] == "loser"
		mu.Unlock()
		if done || time.Now().After(deadline) {
			break
		}
		time.Sleep(2 * time.Millisecond)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(cancelled) != 1 || cancelled[0] != "loser" {
		t.Errorf("cancelled = %v, want [loser]", cancelled)
	}
}

func TestParseStartEvent(t *testing.T) {
	msg := map[string]any{
		"user":     map[string]any{"name": "rival"},
		"quarkid":  "9876543210-5678",
		"playerid": float64(1),
		"port":     float64(6001),
		"delay":    float64(3),
		"ranked":   true,
		"token":    "match-token",
	}

	event := parseStartEvent(msg, "fallback-token")

	if event.Opponent != "rival" {
		t.Errorf("Opponent = %q, want %q", event.Opponent, "rival")
	}
	if event.QuarkID != "9876543210-5678" {
		t.Errorf("QuarkID = %q, want %q", event.QuarkID, "9876543210-5678")
	}
	if event.PlayerID != 1 {
		t.Errorf("PlayerID = %d, want 1", event.PlayerID)
	}
	if event.Port != 6001 {
		t.Errorf("Port = %d, want 6001", event.Port)
	}
	if event.Delay != 3 {
		t.Errorf("Delay = %d, want 3", event.Delay)
	}
	if !event.Ranked {
		t.Error("Ranked = false, want true")
	}
	if event.Token != "match-token" {
		t.Errorf("Token = %q, want %q", event.Token, "match-token")
	}
}

func TestParseStartEvent_FallbackToken(t *testing.T) {
	msg := map[string]any{
		"user":     map[string]any{"name": "rival"},
		"quarkid":  "1234-5678",
		"playerid": float64(0),
		"port":     float64(6000),
		"ranked":   false,
	}

	event := parseStartEvent(msg, "my-session-token")

	if event.Token != "my-session-token" {
		t.Errorf("Token = %q, want fallback %q", event.Token, "my-session-token")
	}
	if event.Delay != 0 {
		t.Errorf("Delay = %d, want 0 (missing field)", event.Delay)
	}
}

func TestSortByRankDistance(t *testing.T) {
	users := []LobbyUser{
		{Name: "me", Rank: 3},
		{Name: "far", Rank: 6},
		{Name: "close", Rank: 4},
		{Name: "closest", Rank: 3},
		{Name: "busy", Rank: 3, Playing: true},
		{Name: "afk", Rank: 3, Away: true},
	}
	sorted := sortByRankDistance(users, 3, "me")

	if len(sorted) != 3 {
		t.Fatalf("expected 3 candidates, got %d", len(sorted))
	}
	if sorted[0].Name != "closest" {
		t.Errorf("expected closest first, got %s", sorted[0].Name)
	}
	if sorted[1].Name != "close" {
		t.Errorf("expected close second, got %s", sorted[1].Name)
	}
	if sorted[2].Name != "far" {
		t.Errorf("expected far third, got %s", sorted[2].Name)
	}
}

func TestParseLobbyUsers(t *testing.T) {
	raw := []any{
		map[string]any{"name": "alice", "rank": float64(4), "playing": map[string]any{"quarkid": "x"}},
		map[string]any{"name": "bob", "rank": float64(2), "away": true},
		map[string]any{"name": "me", "rank": float64(5)},
	}
	users, myRank := parseLobbyUsers(raw, "me")

	if myRank != 5 {
		t.Errorf("myRank = %d, want 5", myRank)
	}
	if len(users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(users))
	}
	if !users[0].Playing {
		t.Error("alice should be playing")
	}
	if !users[1].Away {
		t.Error("bob should be away")
	}
}
