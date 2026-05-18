package database

import (
	"database/sql"
	"fmt"
)

type Game struct {
	ID             int64  `json:"id,string"`
	Title          string `json:"title"`
	Platform       string `json:"platform"`
	ReleaseYear    int    `json:"releaseYear"`
	Developer      string `json:"developer"`
	CoverFilename  string `json:"coverFilename"`
	BannerFilename string `json:"bannerFilename"`
	AppID          string `json:"appId"`
}

func (d *DB) ListGames(query string) ([]Game, error) {
	q := `SELECT id, title, platform, release_year, developer, cover_filename, banner_filename, app_id
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
		if err := rows.Scan(&g.ID, &g.Title, &g.Platform, &g.ReleaseYear, &g.Developer, &g.CoverFilename, &g.BannerFilename, &g.AppID); err != nil {
			return nil, err
		}
		games = append(games, g)
	}
	return games, rows.Err()
}

func (d *DB) GetGame(id int64) (Game, error) {
	var g Game
	err := d.db.QueryRow(
		`SELECT id, title, platform, release_year, developer, cover_filename, banner_filename, app_id
		FROM games WHERE id = ?`, id,
	).Scan(&g.ID, &g.Title, &g.Platform, &g.ReleaseYear, &g.Developer, &g.CoverFilename, &g.BannerFilename, &g.AppID)
	if err == sql.ErrNoRows {
		return g, fmt.Errorf("game %d not found", id)
	}
	return g, err
}

func (d *DB) CreateGame(g Game) (Game, error) {
	res, err := d.db.Exec(
		`INSERT INTO games (title, platform, release_year, developer, cover_filename, banner_filename, app_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		g.Title, g.Platform, g.ReleaseYear, g.Developer, g.CoverFilename, g.BannerFilename, g.AppID,
	)
	if err != nil {
		return g, err
	}
	g.ID, err = res.LastInsertId()
	return g, err
}

func (d *DB) UpdateGame(id int64, g Game) (Game, error) {
	res, err := d.db.Exec(
		`UPDATE games SET title = ?, platform = ?, release_year = ?, developer = ?, cover_filename = ?, banner_filename = ?, app_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		g.Title, g.Platform, g.ReleaseYear, g.Developer, g.CoverFilename, g.BannerFilename, g.AppID, id,
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

func (d *DB) FindGameByAppId(platform string, appId string) (Game, error) {
	var g Game
	err := d.db.QueryRow(
		`SELECT id, title, platform, release_year, developer, cover_filename, banner_filename, app_id
		FROM games WHERE app_id = ? AND platform = ?`, appId, platform,
	).Scan(&g.ID, &g.Title, &g.Platform, &g.ReleaseYear, &g.Developer, &g.CoverFilename, &g.BannerFilename, &g.AppID)

	if err == sql.ErrNoRows {
		return g, fmt.Errorf("game %s not found", appId)
	}
	return g, err
}
