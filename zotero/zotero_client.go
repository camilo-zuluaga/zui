package zotero

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ZoteroClient struct {
	BaseURL string
	UserID  string
	ApiKey  string
	Client  *http.Client
}

func (z *ZoteroClient) FetchAllItems() ([]ZoteroItem, error) {
	url := z.buildURL("items")
	return fetch[ZoteroItem](z, url)
}

func (z *ZoteroClient) FetchItemsByCategory(collectionKey string) ([]ZoteroItem, error) {
	url := z.buildURL("collections", collectionKey, "items")
	return fetch[ZoteroItem](z, url)
}

func (z *ZoteroClient) FetchCollections() ([]Collection, error) {
	url := z.buildURL("collections")
	return fetch[Collection](z, url)
}

func fetch[T any](c *ZoteroClient, url string) ([]T, error) {
	res, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}
	defer res.Body.Close()

	var items []T
	if err := json.NewDecoder(res.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return items, nil
}

func (z *ZoteroClient) buildURL(segments ...string) string {
	u, _ := url.Parse(z.BaseURL)
	pathParts := append([]string{"users", z.UserID}, segments...)
	u.JoinPath(pathParts...)
	return u.String()
}
