package cache

import (
	"database/sql"
	"fmt"

	"github.com/camilo-zuluaga/zotero-tui/zotero"
)

func saveZoteroItems(db *sql.DB, items []zotero.ZoteroItem) error {
	const query = `
INSERT INTO items (
  zotero_key, collection_key, item_type, title, short_title, date, creator_summary, DOI, URL, version
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	ids := make([]int, 0, len(items))

	for _, item := range items {
		res, err := stmt.Exec(
			item.Key,
			"Test",
			item.Data.ItemType,
			item.Data.Title,
			item.Data.ShortTitle,
			item.Data.Date,
			item.Data.CreatorSummary,
			item.Data.DOI,
			item.Data.URL,
			1,
		)
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		ids = append(ids, int(id))
	}

	stmtAttach, err := db.Prepare("INSERT INTO attachments (item_id, zotero_key, title, filename, URL) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare attachment insert: %w", err)
	}
	defer stmtAttach.Close()

	stmtCreator, err := db.Prepare("INSERT INTO creators (item_id, first_name, last_name, creator_type) VALUES (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare creator insert: %w", err)
	}
	defer stmtCreator.Close()

	stmtNote, err := db.Prepare("INSERT INTO notes (zotero_key, item_id, note) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare note insert: %w", err)
	}
	defer stmtNote.Close()

	for i, item := range items {
		itemID := ids[i]
		for _, attachment := range item.Data.Attachment {
			_, err := stmtAttach.Exec(
				itemID,
				attachment.Key,
				attachment.Title,
				attachment.Filename,
				attachment.URL,
			)
			if err != nil {
				return err
			}
		}

		for _, creator := range item.Data.Creators {
			_, err := stmtCreator.Exec(
				itemID,
				creator.FirstName,
				creator.LastName,
				creator.CreatorType,
			)
			if err != nil {
				return err
			}
		}

		for _, note := range item.Data.Note {
			_, err := stmtNote.Exec(
				note.Key,
				itemID,
				note.Note,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getItemsByCollection(db *sql.DB, collectionKey string) ([]zotero.ZoteroItem, error) {
	stmt, err := db.Prepare(`
        SELECT zotero_key, collection_key, item_type, title, short_title, 
               date, creator_summary, DOI, URL 
        FROM items 
        WHERE collection_key = ?
    `)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(collectionKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []zotero.ZoteroItem
	for rows.Next() {
		var r ItemRow
		err := rows.Scan(&r.ZoteroKey, &r.CollectionKey, &r.ItemType, &r.Title,
			&r.ShortTitle, &r.Date, &r.CreatorSummary, &r.DOI, &r.URL)
		if err != nil {
			return nil, err
		}
		items = append(items, r.ToZoteroItem())
	}

	return items, rows.Err()
}
