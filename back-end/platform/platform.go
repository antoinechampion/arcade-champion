package platform

import "back-end/database"

type Platform interface {
	Search(query string) ([]SearchResult, error)
	Launch(appId string) error
}

type SearchResult struct {
	Game  string `json:"game"`
	AppID string `json:"appId"`
}

func Get(platformName string, db *database.DB) Platform {
	switch platformName {
	case "steam":
		return Steam{}
	case "fightcade":
		return Fightcade{db: db}
	case "mame":
		return Mame{db: db}
	}
	return nil
}
