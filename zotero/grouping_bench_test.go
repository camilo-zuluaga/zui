package zotero

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockZoteroServer(gen int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(generateDummyData(gen, 2))
	}))
}

func TestGroupItemsByParentWithMockServer(t *testing.T) {
	want := 1000
	server := newMockZoteroServer(want)
	defer server.Close()

	client := NewZoteroClient(server.URL, "")

	ctx := context.Background()
	items, err := fetch[ZoteroGeneralItem](client, ctx, client.BaseURL)
	if err != nil {
		t.Errorf("couldnt fetch")
	}

	result := GroupItemsByParent(items)

	if len(result) != want {
		t.Errorf("Expected %d parent items, got %d", want, len(result))
	}
}

func BenchmarkGroupItemsByParent(b *testing.B) {
	// We assume 2 childs per parent
	testBenchCases := []struct {
		name       string
		numParents int
	}{
		{"1k Items", 1000},
		{"5000k Items", 5000},
		{"10000k Items", 10000},
	}

	for _, bm := range testBenchCases {
		b.Run(bm.name, func(b *testing.B) {
			server := newMockZoteroServer(bm.numParents)
			defer server.Close()

			client := NewZoteroClient(server.URL, "")

			ctx := context.Background()
			items, err := fetch[ZoteroGeneralItem](client, ctx, client.BaseURL)
			if err != nil {
				b.Errorf("couldnt fetch")
			}

			b.ReportAllocs()

			for b.Loop() {
				_ = GroupItemsByParent(items)
			}

		})
	}
}

func generateDummyData(numParents, childrenPerParent int) []byte {
	// Assume: numParents with 2 children and one note
	exampleParent := `{
		"key": "PARENT%d",
		"data": {
			"itemType": "preprint",
			"title": "TestTitle%d",
			"creators": [
				{"creatorType": "author", "firstName":"authorF%d", "lastName":"authorL%d"}
			],
			"date": "2026-01-01",
			"url": "example.com",
			"collections": ["AAAA"]
		}
	}`

	exampleChild := `{
		"key": "CHILD%d",
		"data": {
			"parentItem": "PARENT%d",
			"itemType": "attachment",
			"title": "Preprint PDF",
			"url": "http://arxiv.org/pdf/example.pdf",
			"filename": "document%d.pdf"
		}
	}`

	exampleNote := `{
		"key": "NOTE%d",
		"data": {
			"parentItem": "PARENT%d",
			"itemType": "note",
			"note": "Comment: This is note %d",
			"tags": [],
			"relations": {}
		}
	}`

	var jsonStr string
	jsonStr = "["

	for i := range numParents {
		parent := fmt.Sprintf(exampleParent, i, i, i, i)
		jsonStr += parent + ","

		for j := range childrenPerParent {
			childNum := i*childrenPerParent + j
			child := fmt.Sprintf(exampleChild, childNum, i, childNum)
			jsonStr += child + ","
		}

		noteNum := i
		note := fmt.Sprintf(exampleNote, noteNum, i, noteNum)
		jsonStr += note

		if i < numParents-1 {
			jsonStr += ","
		}
	}

	jsonStr += "]"
	return []byte(jsonStr)
}
