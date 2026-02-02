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
	Version    int    `json:"version,omitempty"`
	ParentItem string `json:"parentItem"`
	ItemType   string `json:"itemType"`
	Note       string `json:"note"`
}

type ResponseNote struct {
	Data struct {
		Version    int    `json:"version"`
		ParentItem string `json:"parentItem"`
		ItemType   string `json:"itemType"`
		Note       string `json:"note"`
	} `json:"data"`
}

//
// type NoteResponse struct {
// 	Successful map[string]struct {
// 		Data struct {
// 			Key        string `json:"key"`
// 			Version    int    `json:"version"`
// 			ParentItem string `json:"parentItem"`
// 			Note       string `json:"note"`
// 		} `json:"data"`
// 	} `json:"successful"`
// }

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

func (z *ZoteroClient) EditNote(itemKey, newContent string) error {
	url := z.buildURL("items", itemKey)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", z.ApiKey))

	res, err := z.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var result ResponseNote
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	note := Note{
		Version:    result.Data.Version,
		ParentItem: result.Data.ParentItem,
		ItemType:   result.Data.ItemType,
		Note:       newContent,
	}

	marshalled, err := json.Marshal(note)
	if err != nil {
		return err
	}

	req, err = http.NewRequest(http.MethodPatch, url, bytes.NewReader(marshalled))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", z.ApiKey))

	res, err = z.Client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
