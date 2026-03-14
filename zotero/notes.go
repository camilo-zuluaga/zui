package zotero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Note struct {
	Version    int    `json:"version,omitempty"`
	ParentItem string `json:"parentItem"`
	ItemType   string `json:"itemType"`
	Note       string `json:"note"`
}

type APIResponseNote struct {
	Data struct {
		Version    int    `json:"version"`
		ParentItem string `json:"parentItem"`
		ItemType   string `json:"itemType"`
		Note       string `json:"note"`
	} `json:"data"`
}

func (z *ZoteroClient) makeRequest(httpMethod, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", z.ApiKey))

	res, err := z.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request: %w", err)
	}
	return res, nil
}

func marshalJSON(note Note) (io.Reader, error) {
	marshalled, err := json.Marshal([]Note{note})
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(marshalled), nil
}

func (z *ZoteroClient) CreateNote(parentItemKey, content string) (string, error) {
	url, err := buildItemsURL(z.BaseURL, z.UserID, ItemsQuery{})
	if err != nil {
		return "", err
	}

	note := Note{
		ParentItem: parentItemKey,
		ItemType:   "note",
		Note:       content,
	}

	body, err := marshalJSON(note)
	if err != nil {
		return "", err
	}

	res, err := z.makeRequest(http.MethodPost, url, body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("failed to create note, status %d: %s", res.StatusCode, string(responseBody))
	}

	var result struct {
		Success map[string]string `json:"success"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	key, ok := result.Success["0"]
	if !ok {
		return "", fmt.Errorf("no key returned in response")
	}

	return key, nil
}

func (z *ZoteroClient) EditNote(itemKey, newContent string) error {
	url := z.buildURL("items", itemKey)
	res, err := z.makeRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var result APIResponseNote
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

	res, err = z.makeRequest(http.MethodPatch, url, bytes.NewReader(marshalled))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to edit note, status %d: %s", res.StatusCode, string(responseBody))
	}
	return nil
}
