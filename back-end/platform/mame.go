package platform

import (
	"back-end/database"
	"back-end/platform/mame"
	"context"
	"fmt"
	"log"
)

type Mame struct {
	db *database.DB
}

func (m Mame) Search(query string) ([]SearchResult, error) {
	mamePath, err := m.db.MamePath()
	if err != nil {
		return nil, err
	}
	if mamePath == "" {
		return nil, fmt.Errorf("mame path not configured")
	}

	games, err := mame.Search(mamePath, query)
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, len(games))
	for i, g := range games {
		results[i] = SearchResult{Game: g.Name, AppID: g.RomID}
	}
	return results, nil
}

func (m Mame) Launch(_ context.Context, game database.Game, _ LaunchOptions) error {
	mamePath, err := m.db.MamePath()
	if err != nil {
		return err
	}
	if mamePath == "" {
		return fmt.Errorf("mame path not configured")
	}

	log.Printf("[launch-mame] starting title=%q appID=%q command=%q", game.Title, game.AppID, mamePath)
	err = mame.Launch(mamePath, game.AppID)
	if err != nil {
		log.Printf("[launch-mame] failed title=%q appID=%q error=%v", game.Title, game.AppID, err)
		return err
	}
	log.Printf("[launch-mame] started title=%q appID=%q", game.Title, game.AppID)
	return nil
}
