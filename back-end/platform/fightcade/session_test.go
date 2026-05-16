package fightcade

import (
	"testing"
)

func TestBuildMatchURL_ProxyEmulator(t *testing.T) {
	url := buildMatchURL("fbneo", "sfiii3nr1", "1234567890-1234", 0, 6000, 2, true, "tok123")
	want := "fcade://served/fbneo/sfiii3nr1/1234567890-1234.0,6000,2,1"
	if url != want {
		t.Errorf("got %s, want %s", url, want)
	}
}

func TestBuildMatchURL_NativeGGPO(t *testing.T) {
	url := buildMatchURL("flycast", "mvsc", "9999-42", 1, 6001, 3, false, "mytoken")
	want := "fcade://launch/flycast/mvsc/9999-42.1,6001,3,0/mytoken"
	if url != want {
		t.Errorf("got %s, want %s", url, want)
	}
}

func TestBuildPlayURL(t *testing.T) {
	url := buildPlayURL("fbneo", "sfiii3nr1")
	want := "fcade://play/fbneo/sfiii3nr1"
	if url != want {
		t.Errorf("got %s, want %s", url, want)
	}
}

func TestBuildTrainingURL(t *testing.T) {
	url := buildTrainingURL("fbneo", "sfiii3nr1")
	want := "fcade://training/fbneo/sfiii3nr1"
	if url != want {
		t.Errorf("got %s, want %s", url, want)
	}
}

func TestRankName(t *testing.T) {
	cases := []struct {
		rank int
		want string
	}{
		{0, "?"},
		{1, "E"},
		{6, "S"},
		{99, "99"},
	}
	for _, c := range cases {
		if got := RankName(c.rank); got != c.want {
			t.Errorf("RankName(%d) = %s, want %s", c.rank, got, c.want)
		}
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
