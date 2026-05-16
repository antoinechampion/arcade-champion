package handlers

import (
	"back-end/database"
	"encoding/json"
	"net/http"
)

type settings struct {
	FightcadeUsername string `json:"fightcadeUsername"`
	FightcadePassword string `json:"fightcadePassword"`
	FightcadeCookie   string `json:"fightcadeCookie"`
	MamePath          string `json:"mamePath"`
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	password, err := db.FightcadePassword()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cookie, err := db.FightcadeCookie()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mamePath, err := db.MamePath()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings{
		FightcadeUsername: username,
		FightcadePassword: password,
		FightcadeCookie:   cookie,
		MamePath:          mamePath,
	})
}

func putSettings(db *database.DB, w http.ResponseWriter, r *http.Request) {
	var s settings
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.SetFightcadeUsername(s.FightcadeUsername); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := db.SetFightcadePassword(s.FightcadePassword); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := db.SetFightcadeCookie(s.FightcadeCookie); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := db.SetMamePath(s.MamePath); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}
