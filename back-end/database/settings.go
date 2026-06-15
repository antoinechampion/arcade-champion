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

func (d *DB) getEncryptedSetting(key string) (string, error) {
	encrypted, err := d.getSetting(key)
	if err != nil || encrypted == "" {
		return "", err
	}
	return decrypt(encrypted)
}

func (d *DB) setEncryptedSetting(key, value string) error {
	if value == "" {
		return d.setSetting(key, "")
	}
	encrypted, err := encrypt(value)
	if err != nil {
		return err
	}
	return d.setSetting(key, encrypted)
}

func (d *DB) FightcadeCookie() (string, error)  { return d.getEncryptedSetting("fightcade.cookie") }
func (d *DB) FightcadeUsername() (string, error) { return d.getSetting("fightcade.username") }
func (d *DB) FightcadePassword() (string, error) { return d.getEncryptedSetting("fightcade.password") }

func (d *DB) SetFightcadeCookie(v string) error  { return d.setEncryptedSetting("fightcade.cookie", v) }
func (d *DB) SetFightcadeUsername(v string) error { return d.setSetting("fightcade.username", v) }
func (d *DB) SetFightcadePassword(v string) error { return d.setEncryptedSetting("fightcade.password", v) }

func (d *DB) FightcadeMatchDuration() (string, error)    { return d.getSetting("fightcade.matchDuration") }
func (d *DB) SetFightcadeMatchDuration(v string) error   { return d.setSetting("fightcade.matchDuration", v) }

func (d *DB) MamePath() (string, error)    { return d.getSetting("mame.path") }
func (d *DB) SetMamePath(v string) error   { return d.setSetting("mame.path", v) }

func (d *DB) SteamPath() (string, error)  { return d.getSetting("steam.path") }
func (d *DB) SetSteamPath(v string) error { return d.setSetting("steam.path", v) }
