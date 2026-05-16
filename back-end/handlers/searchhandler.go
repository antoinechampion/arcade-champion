package handlers

import (
	"back-end/database"
	"back-end/platform"
	"encoding/json"
	"net/http"
)

func SearchQueryHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		platformString := r.URL.Query().Get("platform")
		query := r.URL.Query().Get("query")

		p := platform.Get(platformString, db)
		if p == nil {
			http.NotFound(w, r)
			return
		}

		results, err := p.Search(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}
