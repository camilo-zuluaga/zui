package zotero

import (
	"fmt"
	"net/url"
	"path"
	"strconv"
)

type ItemsQuery struct {
	CollectionKey string
	ParentKey     string
	Q             string
	QMode         string
	Start         int
	Limit         int
	Sort          string
	Format        string
	Style         string
	Bib           bool
	Children      bool
	Top           bool
	Version       int64
	ToSync        bool
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
		if opts.Top {
			p += "/top"
		}
	} else if opts.Children {
		p = fmt.Sprintf("/users/%s/items/%s/children", url.PathEscape(userID), url.PathEscape(opts.ParentKey))
	} else if opts.Bib {
		p = fmt.Sprintf("/users/%s/items/%s", url.PathEscape(userID), url.PathEscape(opts.ParentKey))
	} else {
		p = fmt.Sprintf("/users/%s/items", url.PathEscape(userID))
		if opts.Top {
			p += "/top"
		}
	}
	u.Path = path.Join(u.Path, p)

	q := url.Values{}
	if opts.Q != "" {
		q.Set("q", opts.Q)
	}
	if opts.QMode != "" {
		q.Set("qmode", opts.QMode)
	}
	if opts.Start >= 0 && !opts.ToSync {
		q.Set("start", strconv.Itoa(opts.Start))
	}
	if opts.Limit > 0 && !opts.ToSync {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Sort != "" {
		q.Set("sort", opts.Sort)
	}
	if opts.Format != "" {
		q.Set("format", opts.Format)
	}
	if opts.Style != "" {
		q.Set("style", opts.Style)
	}
	if opts.ToSync {
		v := fmt.Sprintf("%d", opts.Version)
		q.Set("since", v)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
