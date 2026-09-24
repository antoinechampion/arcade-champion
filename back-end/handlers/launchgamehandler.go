package handlers

import (
	"back-end/database"
	"back-end/platform"
	"encoding/json"
	"log"
	"net/http"
)

type launchGameDto struct {
	Platform      string            `json:"platform"`
	AppID         string            `json:"appId"`
	LaunchOptions map[string]string `json:"launchOptions"`
}

func LaunchGameHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto launchGameDto
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			log.Printf("launch game decode: %v", err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		game, err := db.FindGameByAppId(dto.Platform, dto.AppID)
		if err != nil {
			log.Printf("failed to fetch the game in database: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		log.Printf("[launch] requested platform=%q title=%q appID=%q options=%v", dto.Platform, game.Title, game.AppID, dto.LaunchOptions)
		_ = db.TouchLastPlayed(game.ID)

		p := platform.Get(dto.Platform, db)
		if p == nil {
			http.Error(w, "unknown platform", http.StatusBadRequest)
			return
		}
		err = p.Launch(r.Context(), game, platform.LaunchOptions(dto.LaunchOptions))
		if err != nil {
			log.Printf("[launch] failed platform=%q title=%q appID=%q error=%v", dto.Platform, game.Title, game.AppID, err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		log.Printf("[launch] succeeded platform=%q title=%q appID=%q", dto.Platform, game.Title, game.AppID)
	}
}
