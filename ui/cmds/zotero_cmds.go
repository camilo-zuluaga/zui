package cmds

import (
	"context"
	"os"
	"os/exec"

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

type ZoteroItemsPageMsg struct {
	Items []zotero.ZoteroItem
	Done  bool
	Err   error
}

func StreamCollectionItemsCmd(z *zotero.ZoteroClient, collectionKey string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		ch, errCh := z.StreamItemsByCollection(ctx, collectionKey)
		return StreamStartedMsg{Ch: ch, ErrCh: errCh}
	}
}

func StreamSearchCmd(z *zotero.ZoteroClient, query string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		ch, errCh := z.StreamSearch(ctx, query)
		return StreamStartedMsg{Ch: ch, ErrCh: errCh}
	}
}

type StreamStartedMsg struct {
	Ch    <-chan []zotero.ZoteroGeneralItem
	ErrCh chan error
}

func waitForPageCmd(ch <-chan []zotero.ZoteroGeneralItem, errCh chan error) tea.Cmd {
	return func() tea.Msg {
		items, ok := <-ch
		if !ok {
			var err error
			select {
			case err = <-errCh:
			default:
			}
			return ZoteroItemsPageMsg{Done: true, Err: err}
		}
		return ZoteroItemsPageMsg{
			Items: zotero.MapTopItems(items),
		}
	}
}

func WaitForPageCmd(ch <-chan []zotero.ZoteroGeneralItem, errCh chan error) tea.Cmd {
	return waitForPageCmd(ch, errCh)
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

type ExternalEditorFinishedMsg struct {
	ParentKey string
	Key       string
	Content   string
	New       bool
	Err       error
}

func OpenExternalEditorCmd(editor, parentKey, itemKey, content string, isNew bool) tea.Cmd {
	tmpFile, err := os.CreateTemp("", "zui-note-*.txt")
	if err != nil {
		return func() tea.Msg {
			return ExternalEditorFinishedMsg{Err: err}
		}
	}

	if content != "" {
		tmpFile.WriteString(content)
	}
	tmpFile.Close()

	c := exec.Command(editor, tmpFile.Name())
	return tea.ExecProcess(c, func(err error) tea.Msg {
		defer os.Remove(tmpFile.Name())

		if err != nil {
			return ExternalEditorFinishedMsg{Err: err}
		}

		data, readErr := os.ReadFile(tmpFile.Name())
		if readErr != nil {
			return ExternalEditorFinishedMsg{Err: readErr}
		}

		return ExternalEditorFinishedMsg{
			ParentKey: parentKey,
			Key:       itemKey,
			Content:   string(data),
			New:       isNew,
		}
	})
}
