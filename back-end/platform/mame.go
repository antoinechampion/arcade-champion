package platform

import (
	"back-end/database"
	"back-end/platform/mame"
	"context"
	"fmt"
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

	err = mame.Launch(mamePath, game.AppID)
	if err != nil {
		return err
	}
	return nil
}
