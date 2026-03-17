package zotero

import (
	"net/http"
)

type ZoteroClient struct {
	BaseURL string
	UserID  string
	ApiKey  string
	Client  *http.Client
	Config  *Config
}

func NewZoteroClient(baseURL, userID, apiKey string) *ZoteroClient {
	cfg, _ := LoadConfig()
	return &ZoteroClient{
		BaseURL: baseURL,
		UserID:  userID,
		ApiKey:  apiKey,
		Client:  &http.Client{},
		Config:  cfg,
	}
}

func (z *ZoteroClient) maxItems() int {
	if z.Config != nil && z.Config.MaxItems > 0 {
		return z.Config.MaxItems
	}
	return DefaultMaxItems
}
