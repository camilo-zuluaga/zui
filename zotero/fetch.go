package zotero

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func fetch[T any](c *ZoteroClient, ctx context.Context, url string) ([]T, error) {
	var allItems []T
	currentURL := url

	for currentURL != "" {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, currentURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))

		res, err := c.Client.Do(req)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil, err
			}
			return nil, err
		}

		var items []T
		if err := json.NewDecoder(res.Body).Decode(&items); err != nil {
			res.Body.Close()
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		allItems = append(allItems, items...)
		currentURL = parseNextURL(&res.Header)
		res.Body.Close()
	}

	return allItems, nil
}

func parseNextURL(h *http.Header) string {
	link := h.Get("Link")
	if link == "" {
		return ""
	}
	// the structure of the zotero response is the following:
	// link: <https://api.zotero.org/users/19402717/collections/IXWDFSNI/items?limit=40&start=40>; rel="next", ...
	// so I'm assuming the first < and first > will contain the url for the next set of items
	if strings.Contains(link, `rel="next"`) {
		firstAnchor, lastAnchor := strings.Index(link, "<"), strings.Index(link, ">")
		return link[firstAnchor+1 : lastAnchor]
	}
	return ""
}

func (z *ZoteroClient) buildURL(segments ...string) string {
	u, _ := url.Parse(z.BaseURL)
	pathParts := append([]string{"users", z.UserID}, segments...)
	u = u.JoinPath(pathParts...)
	return u.String()
}
