package zotero

func GroupChildren(zg []ZoteroGeneralItem) ([]ZoteroAttachment, []ZoteroNote) {
	var attachments []ZoteroAttachment
	var notes []ZoteroNote

	for i := range zg {
		item := &zg[i]
		if item.isAttachment() {
			attachments = append(attachments, buildAttachment(item))
		} else if item.isNote() {
			notes = append(notes, buildNote(item))
		}
	}

	return attachments, notes
}

// todo: optimize
func GroupItems(zg []ZoteroGeneralItem) []ZoteroItem {
	parentGroup := make(map[string]*ZoteroItem)
	attachGroup := make(map[string][]ZoteroAttachment)
	notesGroup := make(map[string][]ZoteroNote)

	for i := range zg {
		item := &zg[i]
		if item.Data.ParentItem == "" {
			parentGroup[item.Key] = buildParent(item)
		} else if item.isAttachment() {
			attachGroup[item.Data.ParentItem] = append(
				attachGroup[item.Data.ParentItem],
				buildAttachment(item))
		} else if item.isNote() {
			notesGroup[item.Data.ParentItem] = append(
				notesGroup[item.Data.ParentItem],
				buildNote(item))
		}
	}

	return mergeParentsWithAttachments(parentGroup, attachGroup, notesGroup)
}

func buildParent(z *ZoteroGeneralItem) *ZoteroItem {
	return &ZoteroItem{
		Key: z.Key,
		Data: ZoteroItemData{
			DOI:            z.Data.DOI,
			URL:            z.Data.URL,
			ItemType:       z.Data.ItemType,
			Title:          z.Data.Title,
			ShortTitle:     z.Data.ShortTitle,
			Date:           z.Data.Date,
			CreatorSummary: z.Meta.CreatorSummary,
			Creators:       z.Data.Creators,
			Collections:    z.Data.Collections,
			DateModified:   z.Data.DateModified,
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
	notes map[string][]ZoteroNote,
) []ZoteroItem {
	allItems := make([]ZoteroItem, 0, len(parents))

	for key, parent := range parents {
		parent.Data.Note = notes[key]
		parent.Data.Attachment = attachments[key]
		allItems = append(allItems, *parent)
	}

	return allItems
}
