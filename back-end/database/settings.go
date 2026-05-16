package database

import "database/sql"

func (d *DB) getSetting(key string) (string, error) {
	var value string
	err := d.db.QueryRow(`SELECT value FROM settings WHERE key = ?`, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (d *DB) setSetting(key, value string) error {
	_, err := d.db.Exec(
		`INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		key, value,
	)
	return err
}

func (d *DB) FightcadeCookie() (string, error)   { return d.getSetting("fightcade.cookie") }
func (d *DB) FightcadeUsername() (string, error)  { return d.getSetting("fightcade.username") }
func (d *DB) FightcadePassword() (string, error)  { return d.getSetting("fightcade.password") }

func (d *DB) SetFightcadeCookie(v string) error   { return d.setSetting("fightcade.cookie", v) }
func (d *DB) SetFightcadeUsername(v string) error  { return d.setSetting("fightcade.username", v) }
func (d *DB) SetFightcadePassword(v string) error  { return d.setSetting("fightcade.password", v) }
