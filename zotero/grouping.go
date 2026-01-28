package zotero

// todo: optimize 
func GroupItems(zg []ZoteroGeneralItem) []ZoteroItem {
	parentGroup := make(map[string]*ZoteroItem)
	attachGroup := make(map[string][]ZoteroAttachment)
	notesGroup := make(map[string]ZoteroNote)

	for i := range zg {
		item := &zg[i]
		if item.Data.ParentItem == "" {
			parentGroup[item.Key] = buildParent(item)
		} else if item.isAttachment() {
			attachGroup[item.Data.ParentItem] = append(
				attachGroup[item.Data.ParentItem],
				buildAttachment(item))
		} else if item.isNote() {
			notesGroup[item.Data.ParentItem] = buildNote(item)
		}
	}

	return mergeParentsWithAttachments(parentGroup, attachGroup, notesGroup)
}

func buildParent(z *ZoteroGeneralItem) *ZoteroItem {
	return &ZoteroItem{
		Key: z.Key,
		Data: ZoteroItemData{
			ItemType:    z.Data.ItemType,
			Title:       z.Data.Title,
			Date:        z.Data.Date,
			Creators:    z.Data.Creators,
			Collections: z.Data.Collections,
		},
	}
}

func buildAttachment(z *ZoteroGeneralItem) ZoteroAttachment {
	return ZoteroAttachment{
		Key:      z.Key,
		Title:    z.Data.Title,
		Filename: z.Data.Filename,
		URL:      z.Data.URL,
	}
}

func buildNote(z *ZoteroGeneralItem) ZoteroNote {
	return ZoteroNote{
		Key:  z.Key,
		Note: z.Data.Note,
	}
}

func mergeParentsWithAttachments(parents map[string]*ZoteroItem,
	attachments map[string][]ZoteroAttachment,
	notes map[string]ZoteroNote,
) []ZoteroItem {
	allItems := make([]ZoteroItem, 0, len(parents))

	for key, parent := range parents {
		parent.Data.Note = notes[key]
		parent.Data.Attachment = attachments[key]
		allItems = append(allItems, *parent)
	}

	return allItems
}
