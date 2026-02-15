package cmds

import (
	"context"

	"github.com/camilo-zuluaga/zotero-tui/zotero"
	tea "github.com/charmbracelet/bubbletea"
)

type CollectionLoadedMsg struct {
	Items []zotero.Collection
	Err   error
}

func FetchCollectionsCmd(z *zotero.ZoteroClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		items, err := z.FetchCollections(ctx)
		return CollectionLoadedMsg{Items: items, Err: err}
	}
}

type ZoteroItemsLoadedMsg struct {
	Items []zotero.ZoteroItem
	Err   error
}

func FetchCollectionItemsCmd(z *zotero.ZoteroClient, collectionKey string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		items, err := z.FetchItemsByCollection(ctx, collectionKey)
		return ZoteroItemsLoadedMsg{Items: items, Err: err}
	}
}

type NoteSaved struct {
	Successful bool
}

func SaveNoteCmd(z *zotero.ZoteroClient, parentKey, content string) tea.Cmd {
	return func() tea.Msg {
		err := z.CreateNote(parentKey, content)
		if err != nil {
			return NoteSaved{
				Successful: false,
			}
		}
		return NoteSaved{
			Successful: true,
		}
	}
}

func EditNoteCmd(z *zotero.ZoteroClient, itemKey, newContent string) tea.Cmd {
	return func() tea.Msg {
		err := z.EditNote(itemKey, newContent)
		if err != nil {
			return NoteSaved{
				Successful: false,
			}
		}
		return NoteSaved{
			Successful: true,
		}
	}
}

func FetchQuery(z *zotero.ZoteroClient, query string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		items, err := z.SearchItem(ctx, query)
		return ZoteroItemsLoadedMsg{Items: items, Err: err}
	}
}

func OpenPDF(o *zotero.SystemPDFOpener, key, filename string) error {
	err := o.Open(key, filename)
	if err != nil {
		return err
	}
	return nil
}

type ChildrenLoadedMsg struct {
	ParentKey   string
	Attachments []zotero.ZoteroAttachment
	Notes       []zotero.ZoteroNote
	Err         error
}

func FetchItemChildrenCmd(z *zotero.ZoteroClient, parentKey string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		attachments, notes, err := z.FetchChildren(ctx, parentKey)
		return ChildrenLoadedMsg{
			ParentKey:   parentKey,
			Attachments: attachments,
			Notes:       notes,
			Err:         err,
		}
	}
}

type BibMsg struct {
	Bib string
	Err error
}

func GetBibCmd(z *zotero.ZoteroClient, itemKey, format, style string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		bib, err := z.GetBib(ctx, itemKey, format, style)
		return BibMsg{
			Bib: bib,
			Err: err,
		}
	}
}
