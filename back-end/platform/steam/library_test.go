package steam

import (
	"os"
	"path/filepath"
	"testing"
)

func setupFakeSteam(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	// Create libraryfolders.vdf pointing to two library folders
	steamapps := filepath.Join(root, "steamapps")
	os.MkdirAll(steamapps, 0755)

	lib2 := filepath.Join(root, "games")
	os.MkdirAll(filepath.Join(lib2, "steamapps"), 0755)

	vdf := `"libraryfolders"
{
	"0"
	{
		"path"		"` + root + `"
		"label"		""
	}
	"1"
	{
		"path"		"` + lib2 + `"
		"label"		""
	}
}`
	os.WriteFile(filepath.Join(steamapps, "libraryfolders.vdf"), []byte(vdf), 0644)

	// Library 1: one game
	os.WriteFile(filepath.Join(steamapps, "appmanifest_570.acf"), []byte(`"AppState"
{
	"appid"		"570"
	"name"		"Dota 2"
	"StateFlags"		"4"
}`), 0644)

	// Library 2: two games
	os.WriteFile(filepath.Join(lib2, "steamapps", "appmanifest_730.acf"), []byte(`"AppState"
{
	"appid"		"730"
	"name"		"Counter-Strike 2"
	"StateFlags"		"4"
}`), 0644)

	os.WriteFile(filepath.Join(lib2, "steamapps", "appmanifest_1091500.acf"), []byte(`"AppState"
{
	"appid"		"1091500"
	"name"		"Cyberpunk 2077"
	"StateFlags"		"4"
}`), 0644)

	return root
}

func TestSearch(t *testing.T) {
	root := setupFakeSteam(t)

	tests := []struct {
		query string
		want  int
	}{
		{"", 3},
		{"dota", 1},
		{"counter", 1},
		{"cyber", 1},
		{"nonexistent", 0},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			results, err := Search(root, tt.query)
			if err != nil {
				t.Fatalf("Search(%q): %v", tt.query, err)
			}
			if len(results) != tt.want {
				t.Errorf("Search(%q) returned %d results, want %d", tt.query, len(results), tt.want)
			}
		})
	}
}

func TestParseLibraryFolders(t *testing.T) {
	root := setupFakeSteam(t)

	folders, err := parseLibraryFolders(root)
	if err != nil {
		t.Fatalf("parseLibraryFolders: %v", err)
	}
	if len(folders) != 2 {
		t.Fatalf("got %d folders, want 2", len(folders))
	}
}

func TestParseAppManifest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "appmanifest_570.acf")
	os.WriteFile(path, []byte(`"AppState"
{
	"appid"		"570"
	"name"		"Dota 2"
}`), 0644)

	g, err := parseAppManifest(path)
	if err != nil {
		t.Fatalf("parseAppManifest: %v", err)
	}
	if g.AppID != "570" {
		t.Errorf("AppID = %q, want %q", g.AppID, "570")
	}
	if g.Name != "Dota 2" {
		t.Errorf("Name = %q, want %q", g.Name, "Dota 2")
	}
}
