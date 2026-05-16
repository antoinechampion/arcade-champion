package handlers

import (
	"back-end/platform"
	"encoding/json"
	"fmt"
	"net/http"
)

func SearchQueryHandler(w http.ResponseWriter, r *http.Request) {
	platformString := r.URL.Query().Get("platform")
	query := r.URL.Query().Get("query")

	fmt.Printf("%s - %s\n", platformString, query)

	p := platform.Get(platformString)
	if p == nil {
		http.NotFound(w, r)
		return
	}

	results, err := p.Search(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchResults, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(searchResults)
}
