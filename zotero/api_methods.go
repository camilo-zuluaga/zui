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

func (z *ZoteroClient) FetchItemsByCollection(ctx context.Context, collectionKey string, version int64, toSync bool) ([]ZoteroItem, error) {
	url, _ := buildItemsURL(z.BaseURL, z.UserID, ItemsQuery{
		CollectionKey: collectionKey,
		ToSync:        toSync,
		Version:       version,
		Top:           !toSync,
		Limit:         50,
		Start:         0,
	})
	res, err := fetch[ZoteroGeneralItem](z, ctx, url)
	if err != nil {
		return nil, err
	}
	return MapTopItems(res), nil
}

func (z *ZoteroClient) StreamItemsByCollection(ctx context.Context, collectionKey string) (<-chan []ZoteroGeneralItem, chan error) {
	url, _ := buildItemsURL(z.BaseURL, z.UserID, ItemsQuery{
		CollectionKey: collectionKey,
		Top:           true,
		Limit:         50,
		Start:         0,
	})
	ch := make(chan []ZoteroGeneralItem)
	errCh := make(chan error, 1)
	go func() {
		errCh <- fetchStream(z, ctx, url, ch)
	}()
	return ch, errCh
}

func (z *ZoteroClient) StreamSearch(ctx context.Context, query string) (<-chan []ZoteroGeneralItem, chan error) {
	url, _ := buildItemsURL(z.BaseURL, z.UserID, ItemsQuery{
		Q:     query,
		Top:   true,
		Limit: 50,
		Start: 0,
	})
	ch := make(chan []ZoteroGeneralItem)
	errCh := make(chan error, 1)
	go func() {
		errCh <- fetchStream(z, ctx, url, ch)
	}()
	return ch, errCh
}

func (z *ZoteroClient) FetchCollections(ctx context.Context) ([]Collection, error) {
	url := z.buildURL("collections")
	return fetch[Collection](z, ctx, url)
}

func (z *ZoteroClient) SearchItem(ctx context.Context, query string) ([]ZoteroItem, error) {
	url, _ := buildItemsURL(z.BaseURL, z.UserID, ItemsQuery{
		Q:     query,
		Top:   true,
		Limit: 50,
		Start: 0,
	})
	res, err := fetch[ZoteroGeneralItem](z, ctx, url)
	if err != nil {
		return nil, err
	}
	return MapTopItems(res), nil
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
	})
	_, res, err := rawFetch(z, ctx, url)
	if err != nil {
		return nil, err
	}
	var m map[string]int
	if err := json.Unmarshal(res, &m); err != nil {
		log.Fatal(err)
	}
	return m, nil
}

func (z *ZoteroClient) GetLastModifiedVersion(ctx context.Context, collectionKey string) (string, error) {
	url := z.buildURL("collections")
	header, _, err := rawFetch(z, ctx, url)
	if err != nil {
		return "", err
	}
	if v := header.Get("last-modified-version"); v != "" {
		return v, err
	}
	return "", err
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
