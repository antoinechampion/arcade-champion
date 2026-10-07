package fightcade

import (
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		input    string
		wantName string
		wantArgs []string
	}{
		{"/usr/bin/fightcade", "/usr/bin/fightcade", nil},
		{"flatpak run --command=fcade-quark com.fightcade.Fightcade", "flatpak", []string{"run", "--command=fcade-quark", "com.fightcade.Fightcade"}},
		{"  fightcade  ", "fightcade", nil},
		{"", "", nil},
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
