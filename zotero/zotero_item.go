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
