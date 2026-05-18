package database

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

func ImagesDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	p := filepath.Join(dir, "arcade-champion", "images")
	if err := os.MkdirAll(p, 0755); err != nil {
		return "", err
	}
	return p, nil
}

func generateFilename(gameID int64, kind string) string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%d_%s_%s.jpg", gameID, kind, hex.EncodeToString(b))
}

func (d *DB) SaveImage(gameID int64, kind string, data []byte) (string, error) {
	dir, err := ImagesDir()
	if err != nil {
		return "", err
	}
	filename := generateFilename(gameID, kind)
	if err := os.WriteFile(filepath.Join(dir, filename), data, 0644); err != nil {
		return "", err
	}
	return filename, nil
}

func (d *DB) DeleteImage(filename string) {
	if filename == "" {
		return
	}
	dir, _ := ImagesDir()
	os.Remove(filepath.Join(dir, filename))
}
