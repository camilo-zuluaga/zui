package zotero

import (
	"context"
	"encoding/json"
	"log"
)

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
	return GroupItems(res), nil
}

func (z *ZoteroClient) FetchCollections(ctx context.Context) ([]Collection, error) {
	url := z.buildURL("collections")
	return fetch[Collection](z, ctx, url)
}

func (z *ZoteroClient) SearchItem(ctx context.Context, query string) ([]ZoteroItem, error) {
	url, _ := buildItemsURL(z.BaseURL, z.UserID, ItemsQuery{
		Q:     query,
		Limit: 50,
		Start: 0,
	})
	res, err := fetch[ZoteroGeneralItem](z, ctx, url)
	if err != nil {
		return nil, err
	}
	return GroupItems(res), nil
}

func (z *ZoteroClient) FetchChildren(ctx context.Context, parentKey string) ([]ZoteroAttachment, []ZoteroNote, error) {
	url, _ := buildItemsURL(z.BaseURL, z.UserID, ItemsQuery{
		ParentKey: parentKey,
		Children:  true,
	})
	res, err := fetch[ZoteroGeneralItem](z, ctx, url)
	if err != nil {
		return nil, nil, err
	}
	attachments, notes := GroupChildren(res)
	return attachments, notes, nil
}

func (z *ZoteroClient) FetchItemsVersion(ctx context.Context, collectionKey string) (map[string]int, error) {
	url, _ := buildItemsURL(z.BaseURL, z.UserID, ItemsQuery{
		CollectionKey: collectionKey,
		Version:       true,
	})
	res, err := rawFetch(z, ctx, url)
	if err != nil {
		return nil, err
	}
	var m map[string]int
	if err := json.Unmarshal(res, &m); err != nil {
		log.Fatal(err)
	}
	return m, nil
}

func (z *ZoteroClient) GetBib(ctx context.Context, itemKey, format, style string) (string, error) {
	url, _ := buildItemsURL(z.BaseURL, z.UserID, ItemsQuery{
		ParentKey: itemKey,
		Bib:       true,
		Format:    format,
		Style:     style,
	})

	res, err := simpleFetch(z, ctx, url)
	if err != nil {
		return "", nil
	}

	return res, nil
}
