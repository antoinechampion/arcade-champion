package handlers

import (
	"back-end/database"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func RecentlyPlayedHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		games, err := db.ListRecentlyPlayed(10)
		if err != nil {
			log.Printf("recently played: %v", err)
			http.Error(w, "failed to list recently played", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(games)
	}
}

func ListGamesHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		games, err := db.ListGames(query)
		if err != nil {
			log.Printf("list games: %v", err)
			http.Error(w, "failed to list games", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(games)
	}
}

func GetGameHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		game, err := db.GetGame(id)
		if err != nil {
			http.Error(w, "game not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	}
}

func CreateGameHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Printf("create game parse form: %v", err)
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}

		game := database.Game{
			Title:     r.FormValue("title"),
			Platform:  r.FormValue("platform"),
			Developer: r.FormValue("developer"),
			AppID:     r.FormValue("appId"),
		}
		year, _ := strconv.Atoi(r.FormValue("releaseYear"))
		game.ReleaseYear = year

		created, err := db.CreateGame(game)
		if err != nil {
			log.Printf("create game db: %v", err)
			http.Error(w, "failed to create game", http.StatusInternalServerError)
			return
		}

		coverFilename, err := saveFormImage(db, r, "cover", created.ID)
		if err != nil {
			log.Printf("create game save cover: %v", err)
			_ = db.DeleteGame(created.ID)
			http.Error(w, "failed to save cover image", http.StatusInternalServerError)
			return
		}

		bannerFilename, err := saveFormImage(db, r, "banner", created.ID)
		if err != nil {
			log.Printf("create game save banner: %v", err)
			db.DeleteImage(coverFilename)
			_ = db.DeleteGame(created.ID)
			http.Error(w, "failed to save banner image", http.StatusInternalServerError)
			return
		}

		created.CoverFilename = coverFilename
		created.BannerFilename = bannerFilename
		created, err = db.UpdateGame(created.ID, created)
		if err != nil {
			log.Printf("create game update filenames: %v", err)
			http.Error(w, "failed to update game", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(created)
	}
}

func UpdateGameHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		existing, err := db.GetGame(id)
		if err != nil {
			http.Error(w, "game not found", http.StatusNotFound)
			return
		}

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Printf("update game parse form: %v", err)
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}

		game := database.Game{
			Title:          r.FormValue("title"),
			Platform:       r.FormValue("platform"),
			Developer:      r.FormValue("developer"),
			AppID:          r.FormValue("appId"),
			CoverFilename:  existing.CoverFilename,
			BannerFilename: existing.BannerFilename,
		}
		year, _ := strconv.Atoi(r.FormValue("releaseYear"))
		game.ReleaseYear = year

		if _, _, err := r.FormFile("cover"); err == nil {
			db.DeleteImage(existing.CoverFilename)
			filename, err := saveFormImage(db, r, "cover", id)
			if err != nil {
				log.Printf("update game save cover: %v", err)
				http.Error(w, "failed to save cover image", http.StatusInternalServerError)
				return
			}
			game.CoverFilename = filename
		}

		if _, _, err := r.FormFile("banner"); err == nil {
			db.DeleteImage(existing.BannerFilename)
			filename, err := saveFormImage(db, r, "banner", id)
			if err != nil {
				log.Printf("update game save banner: %v", err)
				http.Error(w, "failed to save banner image", http.StatusInternalServerError)
				return
			}
			game.BannerFilename = filename
		}

		updated, err := db.UpdateGame(id, game)
		if err != nil {
			log.Printf("update game db: %v", err)
			http.Error(w, "failed to update game", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updated)
	}
}

func DeleteGameHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		game, err := db.GetGame(id)
		if err != nil {
			http.Error(w, "game not found", http.StatusNotFound)
			return
		}

		if err := db.DeleteGame(id); err != nil {
			log.Printf("delete game: %v", err)
			http.Error(w, "failed to delete game", http.StatusInternalServerError)
			return
		}

		db.DeleteImage(game.CoverFilename)
		db.DeleteImage(game.BannerFilename)
		w.WriteHeader(http.StatusNoContent)
	}
}

func saveFormImage(db *database.DB, r *http.Request, field string, gameID int64) (string, error) {
	file, header, err := r.FormFile(field)
	if err != nil {
		log.Printf("[image-upload] gameID=%d field=%q no file error=%v", gameID, field, err)
		return "", err
	}
	defer file.Close()

	ct := header.Header.Get("Content-Type")
	log.Printf("[image-upload] gameID=%d field=%q filename=%q contentType=%q declaredSize=%d", gameID, field, header.Filename, ct, header.Size)
	if ct != "" && ct != "image/jpeg" && ct != "image/png" && ct != "image/webp" {
		log.Printf("[image-upload] rejected gameID=%d field=%q reason=unsupported-content-type contentType=%q", gameID, field, ct)
		return "", fmt.Errorf("unsupported image type: %s", ct)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("[image-upload] read failed gameID=%d field=%q error=%v", gameID, field, err)
		return "", err
	}
	log.Printf("[image-upload] read gameID=%d field=%q bytes=%d", gameID, field, len(data))
	upscalerPath, err := db.RealesrganPath()
	if err != nil {
		log.Printf("[image-upload] load upscaler setting failed gameID=%d field=%q error=%v", gameID, field, err)
		return "", err
	}
	log.Printf("[image-upload] upscaler setting gameID=%d field=%q configured=%t path=%q", gameID, field, upscalerPath != "", upscalerPath)
	data, err = upscaleImage(data, field, upscalerPath)
	if err != nil {
		log.Printf("[image-upload] upscale failed gameID=%d field=%q error=%v", gameID, field, err)
		return "", err
	}
	filename, err := db.SaveImage(gameID, field, data)
	if err != nil {
		log.Printf("[image-upload] save failed gameID=%d field=%q bytes=%d error=%v", gameID, field, len(data), err)
		return "", err
	}
	log.Printf("[image-upload] saved gameID=%d field=%q filename=%q bytes=%d", gameID, field, filename, len(data))
	return filename, nil
}
