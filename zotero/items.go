package zotero

func (z *ZoteroClient) FetchAllItems() ([]ZoteroItem, error) {
	url := z.buildURL("items")
	return fetch[ZoteroItem](z, url)
}

func (z *ZoteroClient) FetchItemsByCategory(collectionKey string) ([]ZoteroItem, error) {
	url := z.buildURL("collections", collectionKey, "items")
	return fetch[ZoteroItem](z, url)
}
