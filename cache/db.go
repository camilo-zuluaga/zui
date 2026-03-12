package cache

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Init() (*Cache, error) {
	dbPath, err := dbPath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &Cache{db: db}, nil
}

func dbPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(homeDir, ".config", "zui", "db")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "cache.db"), nil
}

func createTables(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS collections (
		zotero_key TEXT PRIMARY KEY,
		name       TEXT NOT NULL DEFAULT '',
		num_items  INTEGER NOT NULL DEFAULT 0,
		version    INTEGER NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS items (
		zotero_key     TEXT PRIMARY KEY,
		version        INTEGER NOT NULL DEFAULT 0,
		item_type      TEXT NOT NULL DEFAULT '',
		title          TEXT NOT NULL DEFAULT '',
		short_title    TEXT NOT NULL DEFAULT '',
		date           TEXT NOT NULL DEFAULT '',
		creator_summary TEXT NOT NULL DEFAULT '',
		doi            TEXT NOT NULL DEFAULT '',
		url            TEXT NOT NULL DEFAULT '',
		date_modified  TEXT NOT NULL DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS item_collections (
		item_key       TEXT NOT NULL,
		collection_key TEXT NOT NULL,
		PRIMARY KEY (item_key, collection_key),
		FOREIGN KEY (item_key) REFERENCES items(zotero_key) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS creators (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		item_key     TEXT NOT NULL,
		creator_type TEXT NOT NULL DEFAULT '',
		first_name   TEXT NOT NULL DEFAULT '',
		last_name    TEXT NOT NULL DEFAULT '',
		FOREIGN KEY (item_key) REFERENCES items(zotero_key) ON DELETE CASCADE
	);`

	_, err := db.Exec(schema)
	return err
}

func (c *Cache) Close() error {
	return c.db.Close()
}
