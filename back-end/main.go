package main

import (
	"back-end/database"
	"back-end/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	imagesDir, err := database.ImagesDir()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/games/recent", handlers.RecentlyPlayedHandler(db)).
		Methods("GET")
	r.HandleFunc("/api/games", handlers.ListGamesHandler(db)).
		Methods("GET")
	r.HandleFunc("/api/games", handlers.CreateGameHandler(db)).
		Methods("POST")
	r.HandleFunc("/api/games/{id}", handlers.GetGameHandler(db)).
		Methods("GET")
	r.HandleFunc("/api/games/{id}", handlers.UpdateGameHandler(db)).
		Methods("PUT")
	r.HandleFunc("/api/games/{id}", handlers.DeleteGameHandler(db)).
		Methods("DELETE")
	r.HandleFunc("/api/search", handlers.SearchQueryHandler(db)).
		Methods("GET")
	r.HandleFunc("/api/settings", handlers.SettingsHandler(db)).
		Methods("GET", "PUT")
	r.HandleFunc("/api/launch", handlers.LaunchGameHandler(db)).
		Methods("POST")
	r.PathPrefix("/images/").Handler(
		http.StripPrefix("/images/", http.FileServer(http.Dir(imagesDir))),
	)

	log.Fatal(http.ListenAndServe(":8080", cors(r)))
}
