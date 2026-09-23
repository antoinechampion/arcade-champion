package handlers

import (
	"io"
	"net/http"
	"net/url"
)

func ProxyImageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetURL := r.URL.Query().Get("url")
		if targetURL == "" {
			http.Error(w, "missing url parameter", http.StatusBadRequest)
			return
		}

		parsed, err := url.Parse(targetURL)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}

		resp, err := http.Get(targetURL)
		if err != nil {
			http.Error(w, "failed to fetch image", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "failed to fetch image", resp.StatusCode)
			return
		}

		ct := resp.Header.Get("Content-Type")
		if ct == "" {
			ct = "image/jpeg"
		}
		w.Header().Set("Content-Type", ct)
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.WriteHeader(http.StatusOK)

		_, _ = io.Copy(w, resp.Body)
	}
}
