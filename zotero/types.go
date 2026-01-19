package zotero

type ZoteroItem struct {
	Key  string
	Data ZoteroItemData
}

type ZoteroItemData struct {
	Key         string
	ItemType    string
	Title       string
	Collections []string
}

type Collection struct {
	Key  string
	Meta struct {
		NumItems       int
		NumCollections int
	}
	Data struct {
		Name    string
		Version int
	}
}
