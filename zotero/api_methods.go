package zotero

import "context"

func (z *ZoteroClient) FetchAllItems(ctx context.Context) ([]ZoteroItem, error) {
	url := z.buildURL("items")
	return fetch[ZoteroItem](z, ctx, url)
}

func (z *ZoteroClient) FetchItemsByCategory(ctx context.Context, collectionKey string) ([]ZoteroItem, error) {
	url := z.buildURL("collections", collectionKey, "items")
	return fetch[ZoteroItem](z, ctx, url)
}
func (z *ZoteroClient) FetchCollections(ctx context.Context) ([]Collection, error) {
	url := z.buildURL("collections")
	return fetch[Collection](z, ctx, url)
}
