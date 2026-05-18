package platform

import "back-end/database"

type Platform interface {
	Search(query string) ([]SearchResult, error)
	Launch(game database.Game) error
}

type SearchResult struct {
	Game  string `json:"game"`
	AppID string `json:"appId"`
}

func Get(platformName string, db *database.DB) Platform {
	switch platformName {
	case "Steam":
		return Steam{db: db}
	case "Fightcade":
		return Fightcade{db: db}
	case "MAME":
		return Mame{db: db}
	}
	return nil
}
