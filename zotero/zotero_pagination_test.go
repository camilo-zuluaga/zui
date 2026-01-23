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
	tests := []struct {
		testDescription string
		totalItems      int
	}{
		{testDescription: "test pagination with 100 items", totalItems: 100},
		{testDescription: "test pagination with 300 items", totalItems: 300},
		{testDescription: "test pagination with 500 items", totalItems: 500},
		{testDescription: "test pagination with 1000 items", totalItems: 1000},
		{testDescription: "test pagination with 5000 items", totalItems: 5000},
	}
	for _, tt := range tests {
		t.Run(tt.testDescription, func(t *testing.T) {
			totalItems := tt.totalItems
			server := createPaginationMockServer(totalItems)
			client := NewZoteroClient(server.URL, "TEST", "12345")

			want := totalItems
			ctx := context.Background()
			got, err := client.FetchItemsByCollection(ctx, "AAAA")
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if len(got) != want {
				t.Errorf("expected %d elements, got %d", want, len(got))
			}
		})
	}
}

func createPaginationMockServer(totalItems int) *httptest.Server {
	server := httptest.NewServer(nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startStr := r.URL.Query().Get("start")
		start, _ := strconv.Atoi(startStr)

		limitStr := r.URL.Query().Get("limit")
		limit, _ := strconv.Atoi(limitStr)
		end := min(limit + start, totalItems)

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
