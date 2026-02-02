package zotero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Note struct {
	ParentItem string `json:"parentItem"`
	ItemType   string `json:"itemType"`
	Note       string `json:"note"`
}

func (z *ZoteroClient) CreateNote(parentItemKey, content string) error {
	url, err := buildItemsURL(z.BaseURL, z.UserID, ItemsQuery{})
	if err != nil {
		return err
	}

	note := Note{
		ParentItem: parentItemKey,
		ItemType:   "note",
		Note:       content,
	}
	marshalled, err := json.Marshal([]Note{note})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshalled))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", z.ApiKey))

	res, err := z.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("warning: failed to read response body: %v", err)
		body = []byte{}
	}
	fmt.Println(string(body))

	return nil
}
