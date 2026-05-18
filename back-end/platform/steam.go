package platform

import (
	"back-end/database"
	"back-end/platform/steam"
	"fmt"
	"os/exec"
	"runtime"
)

type Steam struct {
	db *database.DB
}

func (s Steam) Search(query string) ([]SearchResult, error) {
	steamPath, err := s.db.SteamPath()
	if err != nil {
		return nil, err
	}
	if steamPath == "" {
		return nil, fmt.Errorf("steam path not configured")
	}

	games, err := steam.Search(steamPath, query)
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, len(games))
	for i, g := range games {
		results[i] = SearchResult{Game: g.Name, AppID: g.AppID}
	}
	return results, nil
}

func (s Steam) Launch(game database.Game) error {
	url := fmt.Sprintf("steam://run/%s", game.AppID)
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "windows":
		return exec.Command("cmd", "/c", "start", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}
