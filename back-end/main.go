package main

import (
	"back-end/database"
	"back-end/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
