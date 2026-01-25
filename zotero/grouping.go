package zotero

// this probably could be further optimized !
func GroupItemsByParent(zg []ZoteroGeneralItem) []ZoteroItem {
	parentGroup := make(map[string]ZoteroItem)
	attachGroup := make(map[string]ZoteroAttachment)

	for i := range zg {
		item := &zg[i]
		if item.Data.ParentItem == "" {
			parentGroup[item.Key] = ZoteroItem{
				Key: zg[i].Key,
				Data: ZoteroItemData{
					ItemType:    zg[i].Data.ItemType,
					Title:       zg[i].Data.Title,
					Date:        zg[i].Data.Date,
					Creators:    zg[i].Data.Creators,
					Collections: zg[i].Data.Collections,
				},
			}
		} else if item.Data.ItemType == "attachment" {
			attachGroup[item.Data.ParentItem] = ZoteroAttachment{
				Key:      zg[i].Data.Title,
				Title:    zg[i].Data.Title,
				Filename: zg[i].Data.Filename,
				URL:      zg[i].Data.URL,
			}
		}
	}

	allItems := make([]ZoteroItem, 0, len(parentGroup))
	for key, val := range parentGroup {
		allItems = append(allItems,
			ZoteroItem{
				Key: key,
				Data: ZoteroItemData{
					ItemType:    val.Data.ItemType,
					Title:       val.Data.Title,
					Date:        val.Data.Date,
					Creators:    val.Data.Creators,
					Attachment:  attachGroup[key],
					Collections: val.Data.Collections,
				},
			},
		)
	}

	return allItems
}
