package zotero

import (
	"encoding/json"
	"fmt"
	"net/url"
)

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
