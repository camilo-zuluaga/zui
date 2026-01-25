package zotero

// General Zotero item trying to deal with heterogeneous data
type ZoteroGeneralItem struct {
	Key  string `json:"key"`
	Data struct {
		ItemType    string              `json:"itemType"`
		ParentItem  string              `json:"parentItem,omitempty"`
		Title       string              `json:"title"`
		URL         string              `json:"url,omitempty"`
		Filename    string              `json:"filename,omitempty"`
		Date        string              `json:"date"`
		Creators    []ZoteroItemCreator `json:"creators,omitempty"`
		Collections []string            `json:"collections,omitempty"`
	}
}

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
	Attachment  ZoteroAttachment
	Collections []string `json:"collections"`
}

type ZoteroItemCreator struct {
	CreatorType string `json:"creatorType"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
}

type ZoteroAttachment struct {
	Key      string
	Title    string
	Filename string
	URL      string
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
	Meta Meta   `json:"meta"`
	Data Data   `json:"data"`
}
