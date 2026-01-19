package zotero

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ZoteroClient struct {
	BaseURL string
	UserID  string
	ApiKey  string
	Client  *http.Client
}

func (z *ZoteroClient) FetchAllItems() ([]ZoteroItem, error) {
	url := fmt.Sprintf("%s/users/%s/items", z.BaseURL, z.UserID)
	return z.fetchItems(url)
}

func (z *ZoteroClient) FetchItemsByCategory(collectionKey string) ([]ZoteroItem, error) {
	url := fmt.Sprintf("%s/users/%s/collections/%s/items", z.BaseURL, z.UserID, collectionKey)
	return z.fetchItems(url)
}

func (z *ZoteroClient) fetchItems(url string) ([]ZoteroItem, error) {
	res, err := z.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}
	defer res.Body.Close()

	var items []ZoteroItem
	if err := json.NewDecoder(res.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return items, nil
}
