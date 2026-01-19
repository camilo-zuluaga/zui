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
	res, err := z.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed")
	}
	defer res.Body.Close()

	var apiResponse []ZoteroItem
	if err := json.NewDecoder(res.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed")
	}

	return apiResponse, nil
}

func (z *ZoteroClient) FetchItemsByCategory(collectionKey string) ([]ZoteroItem, error) {
	url := fmt.Sprintf("%s/users/%s/collections/%s/items", z.BaseURL, z.UserID, collectionKey)

	res, err := z.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed")
	}
	defer res.Body.Close()

	var apiResponse []ZoteroItem
	if err := json.NewDecoder(res.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed")
	}

	return apiResponse, nil
}
