package zotero

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func rawFetch(c *ZoteroClient, ctx context.Context, url string) ([]byte, error) {
	res, err := makeRequest(c, ctx, url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return body, nil
}

func simpleFetch(c *ZoteroClient, ctx context.Context, url string) (string, error) {
	res, err := makeRequest(c, ctx, url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return string(body), nil
}

func fetch[T any](c *ZoteroClient, ctx context.Context, url string) ([]T, error) {
	var allItems []T
	currentURL := url

	for currentURL != "" {
		items, nextURL, err := fetchPage[T](c, ctx, currentURL)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, items...)
		currentURL = nextURL
	}

	return allItems, nil
}

func fetchPage[T any](c *ZoteroClient, ctx context.Context, url string) ([]T, string, error) {
	res, err := makeRequest(c, ctx, url)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	items, err := decodeResponse[T](res)
	if err != nil {
		return nil, "", err
	}

	nextURL := parseNextURL(&res.Header)
	return items, nextURL, nil
}

func makeRequest(c *ZoteroClient, ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	return res, nil
}

func decodeResponse[T any](res *http.Response) ([]T, error) {
	var items []T
	if err := json.NewDecoder(res.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return items, nil
}

func parseNextURL(h *http.Header) string {
	link := h.Get("Link")
	if link == "" {
		return ""
	}

	links := strings.Split(link, ", ")

	m := make(map[string]string, len(links))
	for _, l := range links {
		sep := strings.Split(l, "; ") // format = <link>; rel="action"
		m[sep[1]] = sep[0]
	}

	relNext := m[`rel="next"`]
	if relNext != "" {
		firstAnchor, lastAnchor := strings.Index(relNext, "<"), strings.Index(relNext, ">")
		return relNext[firstAnchor+1 : lastAnchor]
	}
	return ""
}

func (z *ZoteroClient) buildURL(segments ...string) string {
	u, _ := url.Parse(z.BaseURL)
	pathParts := append([]string{"users", z.UserID}, segments...)
	u = u.JoinPath(pathParts...)
	return u.String()
}
