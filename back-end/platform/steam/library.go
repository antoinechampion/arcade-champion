package steam

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Game struct {
	AppID string
	Name  string
}

func Search(steamPath, query string) ([]Game, error) {
	folders, err := parseLibraryFolders(steamPath)
	if err != nil {
		return nil, err
	}

	lowerQuery := strings.ToLower(query)
	var results []Game

	for _, folder := range folders {
		games, err := scanManifests(folder)
		if err != nil {
			continue
		}
		for _, g := range games {
			if strings.Contains(strings.ToLower(g.Name), lowerQuery) {
				results = append(results, g)
			}
		}
	}
	return results, nil
}

func parseLibraryFolders(steamPath string) ([]string, error) {
	vdfPath := filepath.Join(steamPath, "steamapps", "libraryfolders.vdf")
	f, err := os.Open(vdfPath)
	if err != nil {
		return nil, fmt.Errorf("open libraryfolders.vdf: %w", err)
	}
	defer f.Close()

	var folders []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		key, value, ok := parseVDFLine(line)
		if !ok {
			continue
		}
		if key == "path" {
			folders = append(folders, value)
		}
	}
	return folders, scanner.Err()
}

func scanManifests(libraryPath string) ([]Game, error) {
	pattern := filepath.Join(libraryPath, "steamapps", "appmanifest_*.acf")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var games []Game
	for _, path := range matches {
		g, err := parseAppManifest(path)
		if err != nil {
			continue
		}
		games = append(games, g)
	}
	return games, nil
}

func parseAppManifest(path string) (Game, error) {
	f, err := os.Open(path)
	if err != nil {
		return Game{}, err
	}
	defer f.Close()

	var g Game
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		key, value, ok := parseVDFLine(line)
		if !ok {
			continue
		}
		switch key {
		case "appid":
			g.AppID = value
		case "name":
			g.Name = value
		}
		if g.AppID != "" && g.Name != "" {
			break
		}
	}
	if g.AppID == "" || g.Name == "" {
		return Game{}, fmt.Errorf("incomplete manifest: %s", path)
	}
	return g, scanner.Err()
}

// parseVDFLine extracts a key-value pair from a VDF line like: "key" "value"
func parseVDFLine(line string) (key, value string, ok bool) {
	parts := strings.SplitN(line, "\t", 2)
	if len(parts) < 2 {
		parts = strings.SplitN(line, "\"  \"", 2)
		if len(parts) < 2 {
			return "", "", false
		}
	}
	key = strings.Trim(strings.TrimSpace(parts[0]), "\"")
	value = strings.Trim(strings.TrimSpace(parts[1]), "\"")
	if key == "" || value == "" {
		return "", "", false
	}
	return key, value, true
}
