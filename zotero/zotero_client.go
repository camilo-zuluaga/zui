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
