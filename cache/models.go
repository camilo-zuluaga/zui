package cache

import "github.com/camilo-zuluaga/zui/zotero"

func itemToZotero(key string, version int, itemType, title, shortTitle, date, creatorSummary, doi, url, dateModified string,
	collections []string, creators []zotero.ZoteroItemCreator,
) zotero.ZoteroItem {
	return zotero.ZoteroItem{
		Key:     key,
		Version: version,
		Data: zotero.ZoteroItemData{
			DOI:            doi,
			URL:            url,
			ItemType:       itemType,
			Title:          title,
			ShortTitle:     shortTitle,
			Date:           date,
			CreatorSummary: creatorSummary,
			Creators:       creators,
			Collections:    collections,
			DateModified:   dateModified,
		},
	}
}
