package zotero

import "context"

func (z *ZoteroClient) FetchAllItems(ctx context.Context) ([]ZoteroItem, error) {
	url := z.buildURL("items")
	return fetch[ZoteroItem](z, ctx, url)
}

func (z *ZoteroClient) FetchItemsByCollection(ctx context.Context, collectionKey string) ([]ZoteroItem, error) {
	url, _ := buildItemsURL(z.BaseURL, z.UserID, ItemsQuery{
		CollectionKey: collectionKey,
		Limit:         50,
		Start:         0,
	})
	res, err := fetch[ZoteroGeneralItem](z, ctx, url)
	if err != nil {
		return nil, err
	}
	return GroupItemsByParent(res), nil
}

func (z *ZoteroClient) FetchCollections(ctx context.Context) ([]Collection, error) {
	url := z.buildURL("collections")
	return fetch[Collection](z, ctx, url)
}
