package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxyImageHandler(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/test.png" {
			w.Header().Set("Content-Type", "image/png")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("fake-png-bytes"))
			return
		}
		http.NotFound(w, r)
	}))
	defer mockServer.Close()

	handler := ProxyImageHandler()

	t.Run("missing url parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/proxy-image", nil)
		rec := httptest.NewRecorder()
		handler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("invalid url scheme", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/proxy-image?url=ftp://example.com/image.png", nil)
		rec := httptest.NewRecorder()
		handler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("successful proxy", func(t *testing.T) {
		target := mockServer.URL + "/test.png"
		req := httptest.NewRequest("GET", "/api/proxy-image?url="+target, nil)
		rec := httptest.NewRecorder()
		handler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "image/png" {
			t.Fatalf("expected Content-Type image/png, got %q", ct)
		}
		if body := rec.Body.String(); body != "fake-png-bytes" {
			t.Fatalf("expected body fake-png-bytes, got %q", body)
		}
	})

	t.Run("remote not found", func(t *testing.T) {
		target := mockServer.URL + "/notfound.png"
		req := httptest.NewRequest("GET", "/api/proxy-image?url="+target, nil)
		rec := httptest.NewRecorder()
		handler(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", rec.Code)
		}
	})
}
