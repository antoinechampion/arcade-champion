package handlers

import (
	"back-end/database"
	"encoding/json"
	"log"
	"net/http"
)

type settings struct {
	FightcadeUsername      string `json:"fightcadeUsername"`
	FightcadePassword      string `json:"fightcadePassword"`
	FightcadeCookie        string `json:"fightcadeCookie"`
	FightcadeMatchDuration string `json:"fightcadeMatchDuration"`
	MamePath               string `json:"mamePath"`
	SteamPath              string `json:"steamPath"`
}

func SettingsHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getSettings(db, w)
		case http.MethodPut:
			putSettings(db, w, r)
		}
	}
}

func getSettings(db *database.DB, w http.ResponseWriter) {
	username, err := db.FightcadeUsername()
	if err != nil {
		log.Printf("get settings fightcade.username: %v", err)
		http.Error(w, "failed to load settings", http.StatusInternalServerError)
		return
	}
	password, err := db.FightcadePassword()
	if err != nil {
		log.Printf("get settings fightcade.password: %v", err)
		http.Error(w, "failed to load settings", http.StatusInternalServerError)
		return
	}
	cookie, err := db.FightcadeCookie()
	if err != nil {
		log.Printf("get settings fightcade.cookie: %v", err)
		http.Error(w, "failed to load settings", http.StatusInternalServerError)
		return
	}
	matchDuration, err := db.FightcadeMatchDuration()
	if err != nil {
		log.Printf("get settings fightcade.matchDuration: %v", err)
		http.Error(w, "failed to load settings", http.StatusInternalServerError)
		return
	}
	if matchDuration == "" {
		matchDuration = "3"
	}
	mamePath, err := db.MamePath()
	if err != nil {
		log.Printf("get settings mame.path: %v", err)
		http.Error(w, "failed to load settings", http.StatusInternalServerError)
		return
	}
	steamPath, err := db.SteamPath()
	if err != nil {
		log.Printf("get settings steam.path: %v", err)
		http.Error(w, "failed to load settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings{
		FightcadeUsername:      username,
		FightcadePassword:      password,
		FightcadeCookie:        cookie,
		FightcadeMatchDuration: matchDuration,
		MamePath:               mamePath,
		SteamPath:              steamPath,
	})
}

func putSettings(db *database.DB, w http.ResponseWriter, r *http.Request) {
	var s settings
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		log.Printf("put settings decode: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := db.SetFightcadeUsername(s.FightcadeUsername); err != nil {
		log.Printf("put settings fightcade.username: %v", err)
		http.Error(w, "failed to save settings", http.StatusInternalServerError)
		return
	}
	if err := db.SetFightcadePassword(s.FightcadePassword); err != nil {
		log.Printf("put settings fightcade.password: %v", err)
		http.Error(w, "failed to save settings", http.StatusInternalServerError)
		return
	}
	if err := db.SetFightcadeCookie(s.FightcadeCookie); err != nil {
		log.Printf("put settings fightcade.cookie: %v", err)
		http.Error(w, "failed to save settings", http.StatusInternalServerError)
		return
	}
	if err := db.SetFightcadeMatchDuration(s.FightcadeMatchDuration); err != nil {
		log.Printf("put settings fightcade.matchDuration: %v", err)
		http.Error(w, "failed to save settings", http.StatusInternalServerError)
		return
	}
	if err := db.SetMamePath(s.MamePath); err != nil {
		log.Printf("put settings mame.path: %v", err)
		http.Error(w, "failed to save settings", http.StatusInternalServerError)
		return
	}
	if err := db.SetSteamPath(s.SteamPath); err != nil {
		log.Printf("put settings steam.path: %v", err)
		http.Error(w, "failed to save settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}
