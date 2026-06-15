package platform

import (
	"back-end/database"
	"back-end/platform/fightcade"
	"context"
	"fmt"
	"strconv"
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

func (f Fightcade) Launch(ctx context.Context, game database.Game, opts LaunchOptions) error {
	mode := opts["mode"]
	if mode == "" {
		mode = "online"
	}

	creds, err := f.credentials()
	if err != nil {
		return err
	}

	switch mode {
	case "training", "arcade":
		return f.launchOffline(ctx, creds, game.AppID, mode)
	default:
		matchDuration := 3
		if s, err := f.db.FightcadeMatchDuration(); err == nil && s != "" {
			if n, err := strconv.Atoi(s); err == nil {
				matchDuration = n
			}
		}
		_, err = fightcade.Lobby(ctx, creds, game.AppID, matchDuration)
		return err
	}
}

func (f Fightcade) launchOffline(ctx context.Context, creds fightcade.Credentials, appID, mode string) error {
	sr, err := fightcade.Search(ctx, creds, appID)
	if err != nil {
		return err
	}
	f.storeCookie(creds, sr.Cookie)

	var emulator, gameid string
	for _, ch := range sr.Channels {
		if ch.GameID == appID {
			emulator = ch.Emulator
			gameid = ch.GameID
			break
		}
	}
	if emulator == "" {
		return fmt.Errorf("could not resolve emulator for %q", appID)
	}

	if mode == "training" {
		return fightcade.Training(emulator, gameid)
	}
	return fightcade.Play(emulator, gameid)
}
