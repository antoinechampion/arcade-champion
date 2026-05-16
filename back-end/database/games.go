package database

import (
	"database/sql"
	"fmt"
)

type Game struct {
	ID          int64   `json:"id,string"`
	Title       string  `json:"title"`
	Platform    string  `json:"platform"`
	ReleaseYear int     `json:"releaseYear"`
	Developer   string  `json:"developer"`
	ImageURL    string  `json:"imageUrl"`
	BannerURL   *string `json:"bannerUrl,omitempty"`
	PlatformID  string  `json:"platformId"`
}

func (d *DB) ListGames(query string) ([]Game, error) {
	q := `SELECT id, title, platform, release_year, developer, image_url, banner_url, platform_id
		FROM games`
	var args []any
	if query != "" {
		q += ` WHERE title LIKE ?`
		args = append(args, "%"+query+"%")
	}
	q += ` ORDER BY title`

	rows, err := d.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var g Game
		if err := rows.Scan(&g.ID, &g.Title, &g.Platform, &g.ReleaseYear, &g.Developer, &g.ImageURL, &g.BannerURL, &g.PlatformID); err != nil {
			return nil, err
		}
		games = append(games, g)
	}
	return games, rows.Err()
}

func (d *DB) GetGame(id int64) (Game, error) {
	var g Game
	err := d.db.QueryRow(
		`SELECT id, title, platform, release_year, developer, image_url, banner_url, platform_id
		FROM games WHERE id = ?`, id,
	).Scan(&g.ID, &g.Title, &g.Platform, &g.ReleaseYear, &g.Developer, &g.ImageURL, &g.BannerURL, &g.PlatformID)
	if err == sql.ErrNoRows {
		return g, fmt.Errorf("game %d not found", id)
	}
	return g, err
}

func (d *DB) CreateGame(g Game) (Game, error) {
	res, err := d.db.Exec(
		`INSERT INTO games (title, platform, release_year, developer, image_url, banner_url, platform_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		g.Title, g.Platform, g.ReleaseYear, g.Developer, g.ImageURL, g.BannerURL, g.PlatformID,
	)
	if err != nil {
		return g, err
	}
	g.ID, err = res.LastInsertId()
	return g, err
}

func (d *DB) UpdateGame(id int64, g Game) (Game, error) {
	res, err := d.db.Exec(
		`UPDATE games SET title = ?, platform = ?, release_year = ?, developer = ?, image_url = ?, banner_url = ?, platform_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		g.Title, g.Platform, g.ReleaseYear, g.Developer, g.ImageURL, g.BannerURL, g.PlatformID, id,
	)
	if err != nil {
		return g, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return g, err
	}
	if n == 0 {
		return g, fmt.Errorf("game %d not found", id)
	}
	g.ID = id
	return g, nil
}

func (d *DB) DeleteGame(id int64) error {
	res, err := d.db.Exec(`DELETE FROM games WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("game %d not found", id)
	}
	return nil
}
