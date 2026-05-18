package database

import "testing"

func openTestDB(t *testing.T) *DB {
	t.Helper()
	db, err := OpenPath(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestGamesCRUD(t *testing.T) {
	db := openTestDB(t)

	game := Game{
		Title:          "Street Fighter II",
		Platform:       "Fightcade",
		ReleaseYear:    1991,
		Developer:      "Capcom",
		CoverFilename:  "1_cover_abc123.jpg",
		BannerFilename: "1_banner_def456.jpg",
		AppID:          "sfii",
	}

	created, err := db.CreateGame(game)
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == 0 {
		t.Fatal("expected non-zero ID")
	}

	got, err := db.GetGame(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Street Fighter II" || got.AppID != "sfii" || got.BannerFilename != "1_banner_def456.jpg" {
		t.Fatalf("unexpected game: %+v", got)
	}

	got.Title = "Super Street Fighter II"
	updated, err := db.UpdateGame(got.ID, got)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "Super Street Fighter II" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}

	games, err := db.ListGames("")
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 1 {
		t.Fatalf("expected 1 game, got %d", len(games))
	}

	games, err = db.ListGames("super")
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 1 {
		t.Fatalf("expected 1 game matching 'super', got %d", len(games))
	}

	games, err = db.ListGames("tekken")
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 0 {
		t.Fatalf("expected 0 games matching 'tekken', got %d", len(games))
	}

	if err := db.DeleteGame(created.ID); err != nil {
		t.Fatal(err)
	}

	_, err = db.GetGame(created.ID)
	if err == nil {
		t.Fatal("expected error after delete")
	}
}
