package cache

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func Init() (*Cache, error) {
	db, err := sql.Open("sqlite", "test.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &Cache{db: db}, nil
	// rows, err := db.Query("SELECT * FROM items")
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// defer rows.Close()
	//
	// for rows.Next() {
	//
	// 	var id int
	// 	var zotero_key string
	// 	var collection_key string
	// 	var item_type string
	// 	var title string
	// 	var short_title string
	// 	var date string
	// 	var creator_summary string
	// 	var DOI string
	// 	var URL string
	// 	var version int
	//
	// 	err = rows.Scan(&id, &zotero_key, &collection_key, &item_type,
	// 		&title, &short_title, &date, &creator_summary, &DOI, &URL, &version)
	//
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	//
	// 	fmt.Printf("%d %s %s %s %s %s %s %s %s %s %d\n\n", id, zotero_key, collection_key, item_type, title,
	// 		short_title, date, creator_summary, DOI, URL, version)
	// }
}

func createTables(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS items
	(
		id              INTEGER PRIMARY KEY AUTOINCREMENT,
		zotero_key      TEXT,
		collection_key  TEXT,
		item_type       TEXT,
		title           TEXT,
		short_title     TEXT,
		date            TEXT,
		creator_summary TEXT,
		DOI             TEXT,
		URL             TEXT,
		version         INTEGER
	);

	CREATE TABLE IF NOT EXISTS collections
	(
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		zotero_key TEXT,
		name       TEXT,
		num_items  INTEGER,
		version    INTEGER
	);

	CREATE TABLE IF NOT EXISTS attachments
	(
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		item_id    INTEGER,
		zotero_key TEXT,
		title      TEXT,
		filename   TEXT,
		URL        TEXT,
		FOREIGN KEY (item_id) REFERENCES items (id)
	);

	CREATE TABLE IF NOT EXISTS notes
	(
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		zotero_key TEXT,
		item_id    INTEGER,
		note       TEXT,
		FOREIGN KEY (item_id) REFERENCES items (id)
	);

	CREATE TABLE IF NOT EXISTS creators
	(
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		item_id      INTEGER,
		first_name   TEXT,
		last_name    TEXT,
		creator_type TEXT,
		FOREIGN KEY (item_id) REFERENCES items (id)
	);`

	_, err := db.Exec(schema)
	return err
}

func (c *Cache) Close() error {
	return c.db.Close()
}
