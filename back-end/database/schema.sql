CREATE TABLE IF NOT EXISTS games (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    title        TEXT    NOT NULL,
    platform     TEXT    NOT NULL,
    release_year INTEGER NOT NULL,
    developer    TEXT    NOT NULL,
    image_url    TEXT    NOT NULL,
    banner_url   TEXT,
    platform_id  TEXT    NOT NULL,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS settings (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
