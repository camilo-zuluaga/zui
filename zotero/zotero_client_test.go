package zotero

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
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
				ItemType: "book",
				Title:    "BookTest",
			}},
		}

		ctx := context.Background()
		got, err := client.FetchAllItems(ctx)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		assertResponse(t, got, want)
	})

	t.Run("return user's items by collection", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			response := `[
				{
					"key": "ZFJV8EKW",
					"data": {
						"key": "ZFJV8EKW",
						"itemType": "book",
						"title": "Learning Go",
						"date": "2021",
			            "creators": [
							{
								"creatorType": "author",
								"firstName": "Jon",
								"lastName": "Bodner"
							}
						],
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
					ItemType: "book",
					Title:    "Learning Go",
					Date:     "2021",
					Creators: []ZoteroItemCreator{
						{CreatorType: "author", FirstName: "Jon", LastName: "Bodner"},
					},
					Collections: []string{"3UL5E9NK"},
				},
			},
		}

		ctx := context.Background()
		got, err := client.FetchItemsByCollection(ctx, "3UL5E9NK", 0, false)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		assertResponse(t, got, want)
	})
}

func TestZoteroClientFetchCollections(t *testing.T) {
	t.Run("fetch all user collections", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			response := `[
				{
					"key": "3UL5E9NK",
					"meta": {
						"numCollections": 0,
						"numItems": 1
					},
					"data": {
						"version": 3,
						"name": "CollectionTest"
					}
				},
				{
					"key": "SQWT7EXE",
					"meta": {
						"numCollections": 0,
						"numItems": 3
					},
					"data": {
						"version": 10,
						"name": "CollectionTest2"
					}
				}
			]`
			fmt.Fprint(w, response)
		}))
		defer server.Close()

		client := &ZoteroClient{BaseURL: server.URL, UserID: "TESTID", Client: &http.Client{}}

		want := []Collection{
			{
				Key: "3UL5E9NK",
				Meta: Meta{
					NumItems:       1,
					NumCollections: 0,
				},
				Data: Data{
					Name:    "CollectionTest",
					Version: 3,
				},
			},
			{
				Key: "SQWT7EXE",
				Meta: Meta{
					NumItems:       3,
					NumCollections: 0,
				},
				Data: Data{
					Name:    "CollectionTest2",
					Version: 10,
				},
			},
		}

		ctx := context.Background()
		got, err := client.FetchCollections(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertResponse(t, got, want)
	})
}

func TestZoteroClient_Cancelled(t *testing.T) {
	t.Run("cancel request due to timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(30 * time.Millisecond)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, "[]")
		}))
		defer server.Close()

		client := &ZoteroClient{BaseURL: server.URL, UserID: "TESTID", Client: &http.Client{}}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()
		_, err := client.FetchCollections(ctx)
		if err == nil {
			t.Errorf("expected an error due to context timeout but didn't get one")
		}
	})
}

func assertResponse[T any](t testing.TB, got, want []T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
