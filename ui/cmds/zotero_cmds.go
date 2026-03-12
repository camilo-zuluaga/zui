package cmds

import (
	"context"
	"os"
	"os/exec"

	"github.com/camilo-zuluaga/zotero-tui/cache"
	"github.com/camilo-zuluaga/zotero-tui/sync"
	"github.com/camilo-zuluaga/zotero-tui/zotero"
	tea "github.com/charmbracelet/bubbletea"
)

type CollectionLoadedMsg struct {
	Items []zotero.Collection
	Err   error
}

func LoadCollectionsCmd(c *cache.Cache, ss *sync.SyncService) tea.Cmd {
	return func() tea.Msg {
		cached, _ := c.GetCollections()
		if len(cached) > 0 {
			go func() {
				ctx := context.Background()
				ss.SyncCollections(ctx)
			}()
			return CollectionLoadedMsg{Items: cached}
		}

		ctx := context.Background()
		cols, err := ss.SyncCollections(ctx)
		return CollectionLoadedMsg{Items: cols, Err: err}
	}
}

type ZoteroItemsLoadedMsg struct {
	Items []zotero.ZoteroItem
	Err   error
}

func LoadCollectionItemsCmd(c *cache.Cache, z *zotero.ZoteroClient, collectionKey string) tea.Cmd {
	return func() tea.Msg {
		cached, _ := c.GetItemsByCollection(collectionKey)
		if len(cached) > 0 {
			return ZoteroItemsLoadedMsg{Items: cached}
		}

		// Cold cache: stream from API
		ctx := context.Background()
		ch, errCh := z.StreamItemsByCollection(ctx, collectionKey)
		return StreamStartedMsg{Ch: ch, ErrCh: errCh, Cache: c}
	}
}

type ZoteroItemsPageMsg struct {
	Items []zotero.ZoteroItem
	Done  bool
	Err   error
}

type StreamStartedMsg struct {
	Ch    <-chan []zotero.ZoteroGeneralItem
	ErrCh chan error
	Cache *cache.Cache
}

func WaitForPageCmd(ch <-chan []zotero.ZoteroGeneralItem, errCh chan error, c *cache.Cache) tea.Cmd {
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
		mapped := zotero.MapTopItems(items)
		if c != nil {
			c.UpsertItems(mapped)
		}
		return ZoteroItemsPageMsg{Items: mapped}
	}
}

func StreamSearchCmd(z *zotero.ZoteroClient, query string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		ch, errCh := z.StreamSearch(ctx, query)
		return StreamStartedMsg{Ch: ch, ErrCh: errCh}
	}
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

type NoteSaved struct {
	Successful bool
}

func SaveNoteCmd(z *zotero.ZoteroClient, parentKey, content string) tea.Cmd {
	return func() tea.Msg {
		err := z.CreateNote(parentKey, content)
		if err != nil {
			return NoteSaved{Successful: false}
		}
		return NoteSaved{Successful: true}
	}
}

func EditNoteCmd(z *zotero.ZoteroClient, itemKey, newContent string) tea.Cmd {
	return func() tea.Msg {
		err := z.EditNote(itemKey, newContent)
		if err != nil {
			return NoteSaved{Successful: false}
		}
		return NoteSaved{Successful: true}
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
		return BibMsg{Bib: bib, Err: err}
	}
}

func OpenPDF(o *zotero.SystemPDFOpener, key, filename string) error {
	return o.Open(key, filename)
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
