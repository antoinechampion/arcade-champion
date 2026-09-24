package platform

import (
	"back-end/database"
	"back-end/platform/steam"
	"context"
	"fmt"
	"log"
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

func (s Steam) Launch(_ context.Context, game database.Game, _ LaunchOptions) error {
	url := fmt.Sprintf("steam://run/%s", game.AppID)
	log.Printf("[launch-steam] starting title=%q appID=%q url=%q os=%q", game.Title, game.AppID, url, runtime.GOOS)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		log.Printf("[launch-steam] failed title=%q appID=%q error=%v", game.Title, game.AppID, err)
		return err
	}
	log.Printf("[launch-steam] started title=%q appID=%q pid=%d", game.Title, game.AppID, cmd.Process.Pid)
	return nil
}
