CREATE TABLE IF NOT EXISTS games (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    title           TEXT    NOT NULL,
    platform        TEXT    NOT NULL,
    release_year    INTEGER NOT NULL,
    developer       TEXT    NOT NULL,
    cover_filename  TEXT    NOT NULL,
    banner_filename TEXT    NOT NULL,
    app_id          TEXT    NOT NULL,
    last_played_at  DATETIME
);

CREATE TABLE IF NOT EXISTS settings (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
