package database

import "testing"

func TestFightcadeSettings(t *testing.T) {
	db := openTestDB(t)

	cookie, err := db.FightcadeCookie()
	if err != nil {
		t.Fatal(err)
	}
	if cookie != "" {
		t.Fatalf("expected empty cookie, got %q", cookie)
	}

	if err := db.SetFightcadeCookie("abc123"); err != nil {
		t.Fatal(err)
	}
	cookie, err = db.FightcadeCookie()
	if err != nil {
		t.Fatal(err)
	}
	if cookie != "abc123" {
		t.Fatalf("expected 'abc123', got %q", cookie)
	}

	if err := db.SetFightcadeCookie("updated"); err != nil {
		t.Fatal(err)
	}
	cookie, err = db.FightcadeCookie()
	if err != nil {
		t.Fatal(err)
	}
	if cookie != "updated" {
		t.Fatalf("expected 'updated', got %q", cookie)
	}

	if err := db.SetFightcadeUsername("player1"); err != nil {
		t.Fatal(err)
	}
	if err := db.SetFightcadePassword("secret"); err != nil {
		t.Fatal(err)
	}

	username, err := db.FightcadeUsername()
	if err != nil {
		t.Fatal(err)
	}
	if username != "player1" {
		t.Fatalf("expected 'player1', got %q", username)
	}

	password, err := db.FightcadePassword()
	if err != nil {
		t.Fatal(err)
	}
	if password != "secret" {
		t.Fatalf("expected 'secret', got %q", password)
	}
}
