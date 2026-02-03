package zotero

// General Zotero item trying to deal with heterogeneous data
type ZoteroGeneralItem struct {
	Key  string `json:"key"`
	Data struct {
		ItemType    string              `json:"itemType"`
		ParentItem  string              `json:"parentItem,omitempty"`
		Title       string              `json:"title"`
		ShortTitle  string              `json:"shortTitle"`
		DOI         string              `json:"DOI"`
		URL         string              `json:"url,omitempty"`
		Filename    string              `json:"filename,omitempty"`
		Date        string              `json:"date"`
		Note        string              `json:"note"`
		Creators    []ZoteroItemCreator `json:"creators,omitempty"`
		Collections []string            `json:"collections,omitempty"`
	}
}

func (z *ZoteroGeneralItem) isAttachment() bool {
	return z.Data.ItemType == "attachment"
}

func (z *ZoteroGeneralItem) isNote() bool {
	return z.Data.ItemType == "note"
}

type ZoteroItem struct {
	Key  string         `json:"key"`
	Data ZoteroItemData `json:"data"`
}

type ZoteroItemData struct {
	DOI         string
	ItemType    string
	Title       string
	ShortTitle  string
	Date        string
	NumPages    string
	Creators    []ZoteroItemCreator
	Attachment  []ZoteroAttachment
	Note        ZoteroNote
	Collections []string
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

type ZoteroNote struct {
	Key  string
	Note string
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
