package platform

import (
	"back-end/database"
	"back-end/platform/fightcade"
	"context"
	"fmt"
)

type Fightcade struct {
	db *database.DB
}

func (f Fightcade) credentials() (fightcade.Credentials, error) {
	cookie, err := f.db.FightcadeCookie()
	if err != nil {
		return fightcade.Credentials{}, err
	}
	username, err := f.db.FightcadeUsername()
	if err != nil {
		return fightcade.Credentials{}, err
	}
	password, err := f.db.FightcadePassword()
	if err != nil {
		return fightcade.Credentials{}, err
	}
	if cookie == "" && username == "" {
		return fightcade.Credentials{}, fmt.Errorf("fightcade credentials not configured")
	}
	return fightcade.Credentials{Username: username, Password: password, Cookie: cookie}, nil
}

func (f Fightcade) storeCookie(creds fightcade.Credentials, cookie string) {
	if creds.Cookie == "" && cookie != "" {
		_ = f.db.SetFightcadeCookie(cookie)
	}
}

func (f Fightcade) Search(query string) ([]SearchResult, error) {
	creds, err := f.credentials()
	if err != nil {
		return nil, err
	}

	sr, err := fightcade.Search(context.Background(), creds, query)
	if err != nil {
		return nil, err
	}
	f.storeCookie(creds, sr.Cookie)

	channels := sr.Channels
	if len(channels) > 10 {
		channels = channels[:10]
	}

	results := make([]SearchResult, len(channels))
	for i, ch := range channels {
		results[i] = SearchResult{Game: ch.Name, AppID: ch.GameID}
	}
	return results, nil
}

func (f Fightcade) Launch(appId string) error {
	panic("implement me")
}
