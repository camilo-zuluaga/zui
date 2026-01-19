package zotero

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestZoteroClientFetchItems(t *testing.T) {
	t.Run("returns user's items", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `[{"key": "W8GMZJZ3", "data": {"key": "W8GMZJZ3", "itemType": "book", "title": "BookTest"}}]`)
		}))
		defer server.Close()

		client := &ZoteroClient{BaseURL: server.URL, UserID: "TESTID", Client: &http.Client{}}

		want := []ZoteroItem{
			{Key: "W8GMZJZ3", Data: ZoteroItemData{
				Key:      "W8GMZJZ3",
				ItemType: "book",
				Title:    "BookTest",
			}},
		}
		got, err := client.FetchAllItems()
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		assertResponse(t, got, want)
	})

	t.Run("return user's items by category", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			response := `[
				{
					"key": "ZFJV8EKW",
					"data": {
						"key": "ZFJV8EKW",
						"itemType": "book",
						"title": "Learning Go",
						"collections": ["3UL5E9NK"]
					}
				},
				{
					"key": "ABC12345",
					"data": {
						"key": "ABC12345",
						"itemType": "book",
						"title": "Go Programming Patterns",
						"collections": ["3UL5E9NK"]
					}
				}
			]`
			fmt.Fprint(w, response)
		}))
		defer server.Close()

		client := &ZoteroClient{BaseURL: server.URL, UserID: "TESTID", Client: &http.Client{}}

		want := []ZoteroItem{
			{
				Key: "ZFJV8EKW",
				Data: ZoteroItemData{
					Key:         "ZFJV8EKW",
					ItemType:    "book",
					Title:       "Learning Go",
					Collections: []string{"3UL5E9NK"},
				},
			},
			{
				Key: "ABC12345",
				Data: ZoteroItemData{
					Key:         "ABC12345",
					ItemType:    "book",
					Title:       "Go Programming Patterns",
					Collections: []string{"3UL5E9NK"},
				},
			},
		}

		got, err := client.FetchItemsByCategory("3UL5E9NK")
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		assertResponse(t, got, want)
	})
}

func assertResponse(t testing.TB, got, want []ZoteroItem) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %s, want %s", got, want)
	}
}
