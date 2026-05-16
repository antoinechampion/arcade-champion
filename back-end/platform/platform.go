package platform

type Platform interface {
	Search(query string) ([]SearchResult, error)
	Launch(appId string) error
}

type SearchResult struct {
	Game  string `json:"game"`
	AppID string `json:"appId"`
}

func Get(platformName string) Platform {
	switch platformName {
	case "steam":
		return Steam{}
	case "fightcade":
		return Fightcade{}
	case "mame":
		return Mame{}
	}
	return nil
}
