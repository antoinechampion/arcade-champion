package fightcade

import (
	"testing"
)

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
