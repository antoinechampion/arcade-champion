package main

import (
	"back-end/handlers"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/search", handlers.SearchQueryHandler).
		Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
