package cache

import (
	"github.com/camilo-zuluaga/zui/zotero"
)

func (c *Cache) GetCollections() ([]zotero.Collection, error) {
	rows, err := c.db.Query("SELECT zotero_key, name, num_items, version FROM collections ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []zotero.Collection
	for rows.Next() {
		var key, name string
		var numItems, version int
		if err := rows.Scan(&key, &name, &numItems, &version); err != nil {
			return nil, err
		}
		cols = append(cols, zotero.Collection{
			Key:  key,
			Meta: zotero.Meta{NumItems: numItems},
			Data: zotero.Data{Name: name, Version: version},
		})
	}
	return cols, rows.Err()
}

func (c *Cache) UpsertCollections(cols []zotero.Collection) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO collections (zotero_key, name, num_items, version)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(zotero_key) DO UPDATE SET
			name = excluded.name,
			num_items = excluded.num_items,
			version = excluded.version`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, col := range cols {
		if _, err := stmt.Exec(col.Key, col.Data.Name, col.Meta.NumItems, col.Data.Version); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (c *Cache) GetItemsByCollection(collectionKey string) ([]zotero.ZoteroItem, error) {
	rows, err := c.db.Query(`
		SELECT i.zotero_key, i.version, i.item_type, i.title, i.short_title,
		       COALESCE(NULLIF(i.date, ''), 'No Date') as date, i.creator_summary, COALESCE(NULLIF(i.doi, ''), 'No DOI') as doi, i.url, i.date_modified
		FROM items i
		JOIN item_collections ic ON ic.item_key = i.zotero_key
		WHERE ic.collection_key = ?
		ORDER BY i.date_modified DESC`, collectionKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []zotero.ZoteroItem
	for rows.Next() {
		var key, itemType, title, shortTitle, date, creatorSummary, doi, url, dateMod string
		var version int
		if err := rows.Scan(&key, &version, &itemType, &title, &shortTitle,
			&date, &creatorSummary, &doi, &url, &dateMod); err != nil {
			return nil, err
		}

		creators, _ := c.getCreators(key)
		collections, _ := c.getItemCollections(key)

		items = append(items, itemToZotero(
			key, version, itemType, title, shortTitle, date,
			creatorSummary, doi, url, dateMod,
			collections, creators,
		))
	}
	return items, rows.Err()
}

func (c *Cache) UpsertItems(items []zotero.ZoteroItem) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	itemStmt, err := tx.Prepare(`
		INSERT INTO items (zotero_key, version, item_type, title, short_title,
		                   date, creator_summary, doi, url, date_modified)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(zotero_key) DO UPDATE SET
			version = excluded.version,
			item_type = excluded.item_type,
			title = excluded.title,
			short_title = excluded.short_title,
			date = excluded.date,
			creator_summary = excluded.creator_summary,
			doi = excluded.doi,
			url = excluded.url,
			date_modified = excluded.date_modified`)
	if err != nil {
		return err
	}
	defer itemStmt.Close()

	icStmt, err := tx.Prepare(`
		INSERT INTO item_collections (item_key, collection_key)
		VALUES (?, ?)
		ON CONFLICT DO NOTHING`)
	if err != nil {
		return err
	}
	defer icStmt.Close()

	creatorDel, err := tx.Prepare("DELETE FROM creators WHERE item_key = ?")
	if err != nil {
		return err
	}
	defer creatorDel.Close()

	creatorStmt, err := tx.Prepare(
		"INSERT INTO creators (item_key, creator_type, first_name, last_name) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer creatorStmt.Close()

	for _, item := range items {
		d := item.Data
		if _, err := itemStmt.Exec(
			item.Key, item.Version, d.ItemType, d.Title, d.ShortTitle,
			d.Date, d.CreatorSummary, d.DOI, d.URL, d.DateModified,
		); err != nil {
			return err
		}

		for _, colKey := range d.Collections {
			if _, err := icStmt.Exec(item.Key, colKey); err != nil {
				return err
			}
		}

		if _, err := creatorDel.Exec(item.Key); err != nil {
			return err
		}
		for _, cr := range d.Creators {
			if _, err := creatorStmt.Exec(item.Key, cr.CreatorType, cr.FirstName, cr.LastName); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

func (c *Cache) ClearCollections() error {
	_, err := c.db.Exec("DELETE FROM collections")
	return err
}

func (c *Cache) ClearItemsByCollection(collectionKey string) error {
	_, err := c.db.Exec(`
		DELETE FROM items WHERE zotero_key IN (
			SELECT item_key FROM item_collections WHERE collection_key = ?
		)`, collectionKey)
	return err
}

func (c *Cache) getCreators(itemKey string) ([]zotero.ZoteroItemCreator, error) {
	rows, err := c.db.Query(
		"SELECT creator_type, first_name, last_name FROM creators WHERE item_key = ?", itemKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creators []zotero.ZoteroItemCreator
	for rows.Next() {
		var ct, fn, ln string
		if err := rows.Scan(&ct, &fn, &ln); err != nil {
			return nil, err
		}
		creators = append(creators, zotero.ZoteroItemCreator{
			CreatorType: ct, FirstName: fn, LastName: ln,
		})
	}
	return creators, rows.Err()
}

func (c *Cache) getItemCollections(itemKey string) ([]string, error) {
	rows, err := c.db.Query(
		"SELECT collection_key FROM item_collections WHERE item_key = ?", itemKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		cols = append(cols, key)
	}
	return cols, rows.Err()
}
