package zotero

import (
	"fmt"
	"net/url"
	"path"
	"strconv"
)

// var baseURL = "https://api.zotero.org"

type ItemsQuery struct {
	CollectionKey string
	Q             string
	QMode         string
	Start         int
	Limit         int
}

func buildItemsURL(baseURL, userID string, opts ItemsQuery) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	var p string
	if opts.CollectionKey != "" {
		p = fmt.Sprintf("/users/%s/collections/%s/items",
			url.PathEscape(userID), url.PathEscape(opts.CollectionKey))
	} else {
		p = fmt.Sprintf("/users/%s/items", url.PathEscape(userID))
	}
	u.Path = path.Join(u.Path, p)

	q := url.Values{}
	if opts.Q != "" {
		q.Set("q", opts.Q)
	}
	if opts.QMode != "" {
		q.Set("qmode", opts.QMode)
	}
	if opts.Start >= 0 {
		q.Set("start", strconv.Itoa(opts.Start))
	}
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
