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

	r := mux.NewRouter()
	r.HandleFunc("/api/search", handlers.SearchQueryHandler).
		Methods("GET")
	r.HandleFunc("/api/settings", handlers.SettingsHandler(db)).
		Methods("GET", "PUT")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
