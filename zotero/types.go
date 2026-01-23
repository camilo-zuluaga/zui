package zotero

type ZoteroItem struct {
	Key  string         `json:"key"`
	Data ZoteroItemData `json:"data"`
}

type ZoteroItemData struct {
	ItemType    string              `json:"itemType"`
	Title       string              `json:"title"`
	Date        string              `json:"date"`
	NumPages    string              `json:"numPages"`
	Creators    []ZoteroItemCreator `json:"creators"`
	Collections []string            `json:"collections"`
}

type ZoteroItemCreator struct {
	CreatorType string `json:"creatorType"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
}

type Meta struct {
	NumItems       int `json:"numItems,omitempty"`
	NumCollections int `json:"numCollections,omitempty"`
}

type Data struct {
	Name    string `json:"name,omitempty"`
	Version int    `json:"version,omitempty"`
}

type Collection struct {
	Key  string `json:"key"`
	Meta Meta
	Data Data
}
