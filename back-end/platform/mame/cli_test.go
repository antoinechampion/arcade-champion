package mame

import (
	"testing"
)

var sampleOutput = `Name:             Description:
sf                "Street Fighter (US, set 1)"
sf2               "Street Fighter II: The World Warrior (World 910522)"
sf2ce             "Street Fighter II': Champion Edition (World 920513)"
sfiii3            "Street Fighter III 3rd Strike: Fight for the Future (Europe 990608)"
garou             "Garou: Mark of the Wolves (NGM-2530)"
outrun            "Out Run (set 1)"
outrunb           "Out Run (bootleg)"
outrunners        "OutRunners"
`

func TestParseListFull(t *testing.T) {
	games := parseListFull([]byte(sampleOutput), "street fighter")
	if len(games) != 4 {
		t.Fatalf("expected 4 results, got %d", len(games))
	}
	if games[0].RomID != "sf" || games[0].Name != "Street Fighter (US, set 1)" {
		t.Errorf("unexpected first result: %+v", games[0])
	}
}

func TestParseListFullCaseInsensitive(t *testing.T) {
	games := parseListFull([]byte(sampleOutput), "GAROU")
	if len(games) != 1 {
		t.Fatalf("expected 1 result, got %d", len(games))
	}
	if games[0].RomID != "garou" {
		t.Errorf("expected garou, got %s", games[0].RomID)
	}
}

func TestParseListFullSpaceNormalization(t *testing.T) {
	// "outrun" should match "Out Run" (space in name) and "OutRunners"
	games := parseListFull([]byte(sampleOutput), "outrun")
	if len(games) != 3 {
		t.Fatalf("expected 3 results, got %d: %v", len(games), games)
	}
}

func TestParseListFullNoMatch(t *testing.T) {
	games := parseListFull([]byte(sampleOutput), "tekken")
	if len(games) != 0 {
		t.Fatalf("expected 0 results, got %d", len(games))
	}
}

func TestParseLineSkipsHeader(t *testing.T) {
	_, _, ok := parseLine("Name:             Description:")
	if ok {
		t.Error("expected header line to be skipped")
	}
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		input    string
		wantName string
		wantArgs []string
	}{
		{"/usr/bin/mame", "/usr/bin/mame", nil},
		{"flatpak run org.mamedev.MAME", "flatpak", []string{"run", "org.mamedev.MAME"}},
		{"  mame  ", "mame", nil},
	}
	for _, tt := range tests {
		name, args := parseCommand(tt.input)
		if name != tt.wantName {
			t.Errorf("parseCommand(%q) name = %q, want %q", tt.input, name, tt.wantName)
		}
		if len(args) != len(tt.wantArgs) {
			t.Errorf("parseCommand(%q) args = %v, want %v", tt.input, args, tt.wantArgs)
			continue
		}
		for i := range args {
			if args[i] != tt.wantArgs[i] {
				t.Errorf("parseCommand(%q) args[%d] = %q, want %q", tt.input, i, args[i], tt.wantArgs[i])
			}
		}
	}
}
