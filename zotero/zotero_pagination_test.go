package zotero

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestFollowURLPagination(t *testing.T) {
	t.Run("test pagination", func(t *testing.T) {
		server := createPaginationMockServer(80)
		client := NewZoteroClient(server.URL, "TEST", "12345")

		want := 80
		ctx := context.Background()
		got, _ := client.FetchItemsByCollection(ctx, "AAAA")

		if len(got) != want {
			t.Errorf("expected %d elements, got %d", want, len(got))
		}
	})
}

func createPaginationMockServer(totalItems int) *httptest.Server {
	server := httptest.NewServer(nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startStr := r.URL.Query().Get("start")
		start, _ := strconv.Atoi(startStr)

		limitStr := r.URL.Query().Get("limit")
		limit, _ := strconv.Atoi(limitStr)
		end := limit + start
		if end > totalItems {
			end = totalItems
		}

		var items []ZoteroItem
		for i := start; i < end; i++ {
			iStr := strconv.Itoa(i)
			items = append(items, ZoteroItem{
				Key: iStr,
				Data: ZoteroItemData{
					ItemType:    "book",
					Title:       "Item " + iStr,
					Date:        "2026",
					NumPages:    iStr,
					Creators:    []ZoteroItemCreator{},
					Collections: []string{},
				},
			})
		}

		var nextURL string
		if end < totalItems {
			nextURL = fmt.Sprintf("%s%s?limit=%d&start=%d", server.URL, r.URL.Path, limit, end)
			w.Header().Set("Link", fmt.Sprintf(`<%s>; rel="next"`, nextURL))
		} else {
			nextURL = fmt.Sprintf("%s%s?limit=%d", server.URL, r.URL.Path, limit)
			w.Header().Set("Link", fmt.Sprintf(`<%s>; rel="first"`, nextURL))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	})

	server.Config.Handler = handler
	return server
}
