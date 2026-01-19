package zotero

import (
	"net/http"
)

type ZoteroClient struct {
	BaseURL string
	UserID  string
	ApiKey  string
	Client  *http.Client
}

func NewZoteroClient(baseURL, userID string) *ZoteroClient {
	return &ZoteroClient{
		BaseURL: baseURL,
		UserID:  userID,
		Client:  &http.Client{},
	}
}
