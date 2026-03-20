package ui

import (
	"fmt"

	"github.com/99designs/keyring"
	"github.com/camilo-zuluaga/zui/cache"
	"github.com/camilo-zuluaga/zui/clipboard"
	"github.com/camilo-zuluaga/zui/sync"
	"github.com/camilo-zuluaga/zui/ui/attachpicker"
	"github.com/camilo-zuluaga/zui/ui/cmds"
	"github.com/camilo-zuluaga/zui/ui/collections"
	"github.com/camilo-zuluaga/zui/ui/initial"
	"github.com/camilo-zuluaga/zui/ui/items"
	noteeditor "github.com/camilo-zuluaga/zui/ui/note-editor"
	"github.com/camilo-zuluaga/zui/ui/notepicker"
	"github.com/camilo-zuluaga/zui/ui/search"
	"github.com/camilo-zuluaga/zui/zotero"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type currentView int

const (
	ItemsView currentView = iota
	CollectionsView
	NoteEditorView
	NotePickerView
	SearchView
	AttachmentView
	InitialView
)

type rootModel struct {
	width  int
	height int

	zotero          *zotero.ZoteroClient
	systemPDFOpener *zotero.SystemPDFOpener
	cache           *cache.Cache
	sync            *sync.SyncService

	initial      initial.Model
	collections  collections.Model
	zoteroItems  items.Model
	noteEditor   noteeditor.Model
	notepicker   notepicker.Model
	searchInput  search.Model
	attachReader attachpicker.Model

	currentView          currentView
	currentCollectionKey string

	loading   bool
	streaming bool
	streamCh  <-chan []zotero.ZoteroGeneralItem
	streamErr chan error
	spinner   spinner.Model
}

func NewRootModel(z *zotero.ZoteroClient, c *cache.Cache, ss *sync.SyncService) rootModel {
	s := spinner.New()
	o := zotero.NewSystemPDFOpener()
	return rootModel{
		zotero:          z,
		systemPDFOpener: o,
		cache:           c,
		sync:            ss,
		collections:     collections.New(),
		zoteroItems:     items.New(),
		currentView:     CollectionsView,
		noteEditor:      noteeditor.InitialModel("", "", "", false),
		notepicker:      notepicker.New(""),
		loading:         true,
		spinner:         s,
	}
}

func NewInitialRootModel(c *cache.Cache) rootModel {
	s := spinner.New()
	return rootModel{
		cache:       c,
		collections: collections.New(),
		zoteroItems: items.New(),
		currentView: InitialView,
		initial:     initial.InitialModel(),
		noteEditor:  noteeditor.InitialModel("", "", "", false),
		notepicker:  notepicker.New(""),
		spinner:     s,
	}
}

func (m rootModel) Init() tea.Cmd {
	if m.currentView == InitialView {
		return textinput.Blink
	}
	return tea.Batch(m.spinner.Tick,
		cmds.LoadCollectionsCmd(m.cache, m.sync))
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

		headerHeight := 2
		bodyHeight := m.height - headerHeight

		m.collections.SetSize(m.width, bodyHeight)
		m.zoteroItems.SetSize(m.width, bodyHeight)
		m.noteEditor.SetSize(m.width, bodyHeight)
		m.notepicker.SetSize(m.width, bodyHeight)

	case tea.KeyMsg:
		if m.isFiltering() {
			break
		}
		if m.isFilterApplied() {
			key := msg.String()
			if key != "enter" && key != "n" && key != "r" && key != "b" && key != "q" && key != "ctrl+r" {
				// navigation keys (esc, arrows, etc.) go to the list to manage the filter
				break
			}
		}
		switch msg.String() {
		case "esc":
			if m.currentView == ItemsView {
				m.currentView = CollectionsView
				return m, nil
			}
			if m.currentView == SearchView {
				m.currentView = CollectionsView
				return m, nil
			}
		case "q":
			return m, tea.Quit
		case "ctrl+r":
			if m.currentView == CollectionsView {
				m.loading = true
				m.collections = collections.New()
				m.collections.SetSize(m.width, m.height-2)
				_ = m.cache.ClearCollections()
				return m, tea.Batch(m.spinner.Tick,
					cmds.LoadCollectionsCmd(m.cache, m.sync))
			}
			if m.currentView == ItemsView && m.currentCollectionKey != "" {
				m.loading = true
				m.zoteroItems.ClearItems()
				_ = m.cache.ClearItemsByCollection(m.currentCollectionKey)
				return m, tea.Batch(m.spinner.Tick,
					cmds.LoadCollectionItemsCmd(m.cache, m.zotero, m.currentCollectionKey))
			}
		case "enter":
			if m.currentView == CollectionsView {
				if sel := m.collections.SelectedCollection(); sel != nil {
					m.loading = true
					m.currentView = ItemsView
					m.currentCollectionKey = sel.Key
					m.zoteroItems.ClearItems()
					m.zoteroItems.HelpText(items.ModeNormal)
					return m, tea.Batch(m.spinner.Tick,
						cmds.LoadCollectionItemsCmd(m.cache, m.zotero, sel.Key))
				}
			}

			if m.currentView == ItemsView {
				if sel := m.zoteroItems.SelectedZoteroItem(); sel != nil {
					if len(sel.Data.Attachment) == 0 && len(sel.Data.Note) == 0 {
						return m, cmds.FetchItemChildrenCmd(m.zotero, sel.Key)
					}
				}
			}
		case "n":
			if m.currentView == ItemsView {
				if sel := m.zoteroItems.SelectedZoteroItem(); sel != nil {
					if len(sel.Data.Note) == 0 && len(sel.Data.Attachment) == 0 {
						return m, cmds.FetchItemChildrenCmd(m.zotero, sel.Key)
					}
					m.notepicker = notepicker.New(sel.Key)
					m.notepicker.SetSize(m.width, m.height)
					if len(sel.Data.Note) != 0 {
						m.notepicker.SetZoteroNotes(sel.Data.Note)
						m.currentView = NotePickerView
						return m, nil
					}
					if editor := m.externalEditor(); editor != "" {
						return m, cmds.OpenExternalEditorCmd(editor, sel.Key, "", "", true)
					}
					m.noteEditor = noteeditor.InitialModel(sel.Key, "", "", true)
					m.noteEditor.SetSize(m.width, m.height)
					m.currentView = NoteEditorView
					return m, nil
				}
			}
		case "r":
			if m.currentView == ItemsView {
				if sel := m.zoteroItems.SelectedZoteroItem(); sel != nil {
					if len(sel.Data.Attachment) == 0 {
						return m, cmds.FetchItemChildrenCmd(m.zotero, sel.Key)
					}
				}
				if sel := m.zoteroItems.SelectedZoteroItem(); sel != nil && len(sel.Data.Attachment) != 0 {
					m.attachReader = attachpicker.New(sel.Data.Title)
					m.attachReader.SetSize(m.width, m.height)
					if len(sel.Data.Attachment) > 1 {
						m.attachReader.SetZoteroAttachments(sel.Data.Attachment)
						m.currentView = AttachmentView
						return m, nil
					}
					onlyAttachment := sel.Data.Attachment[0]
					cmds.OpenPDF(m.systemPDFOpener, onlyAttachment.Key, onlyAttachment.Filename)
					return m, nil
				}
			}
		case "b":
			if m.currentView == ItemsView {
				if sel := m.zoteroItems.SelectedZoteroItem(); sel != nil {
					return m, cmds.GetBibCmd(m.zotero, sel.Key, m.zotero.Config.Format, m.zotero.Config.Style)
				}
			}
		case "s":
			if m.currentView == CollectionsView {
				m.currentView = SearchView
				m.zoteroItems.HelpText(items.ModeNormal)
				m.searchInput = search.InitialModel()
				return m, nil
			}
		}

	case notepicker.NoteSelectedMsg:
		if editor := m.externalEditor(); editor != "" {
			return m, cmds.OpenExternalEditorCmd(editor, msg.ParentKey, msg.ItemKey, msg.Content, msg.New)
		}
		m.noteEditor = noteeditor.InitialModel(msg.ParentKey, msg.ItemKey, msg.Content, msg.New)
		m.noteEditor.SetSize(m.width, m.height)
		m.currentView = NoteEditorView
		return m, nil

	case noteeditor.CancelNoteMsg:
		m.currentView = ItemsView
		return m, nil

	case noteeditor.SavedNoteMsg:
		m.currentView = ItemsView
		if !msg.New {
			return m, cmds.EditNoteCmd(m.zotero, msg.ParentKey, msg.Key, msg.Content)
		}
		return m, cmds.SaveNoteCmd(m.zotero, msg.ParentKey, msg.Content)

	case cmds.ExternalEditorFinishedMsg:
		m.currentView = ItemsView
		if msg.Err != nil || msg.Content == "" {
			return m, nil
		}
		if !msg.New {
			return m, cmds.EditNoteCmd(m.zotero, msg.ParentKey, msg.Key, msg.Content)
		}
		return m, cmds.SaveNoteCmd(m.zotero, msg.ParentKey, msg.Content)

	case cmds.NoteSaved:
		if !msg.Successful {
			return m, nil
		}
		if msg.Edited {
			m.zoteroItems.UpdateNote(msg.ParentKey, msg.NoteKey, msg.Content)
			return m, nil
		}
		m.zoteroItems.AppendNote(msg.ParentKey, msg.NoteKey, msg.Content)
		return m, nil

	case initial.CredentialsMsg:
		kr, err := keyring.Open(keyring.Config{
			ServiceName: "zotero-tui",
		})
		if err == nil {
			_ = kr.Set(keyring.Item{Key: "api-key", Data: []byte(msg.APIKey)})
			_ = kr.Set(keyring.Item{Key: "user-id", Data: []byte(msg.UserID)})
		}
		m.zotero = zotero.NewZoteroClient("https://api.zotero.org", msg.UserID, msg.APIKey)
		m.systemPDFOpener = zotero.NewSystemPDFOpener()
		m.sync = sync.New(m.cache, m.zotero)
		m.currentView = CollectionsView
		m.loading = true
		return m, tea.Batch(m.spinner.Tick,
			cmds.LoadCollectionsCmd(m.cache, m.sync))

	case cmds.CollectionLoadedMsg:
		m.loading = false
		m.collections.SetZoteroCollections(msg.Items)
		return m, nil

	case cmds.StreamStartedMsg:
		m.streamCh = msg.Ch
		m.streamErr = msg.ErrCh
		m.streaming = true
		return m, cmds.WaitForPageCmd(m.streamCh, m.streamErr, msg.Cache)

	case cmds.ZoteroItemsLoadedMsg:
		m.loading = false
		m.zoteroItems.SetZoteroItems(msg.Items)
		return m, m.zoteroItems.HelpText(items.ModeNormal)

	case cmds.ZoteroItemsPageMsg:
		if msg.Err != nil {
			m.loading = false
			m.streaming = false
			return m, nil
		}
		if msg.Done {
			m.loading = false
			m.streaming = false
			m.streamCh = nil
			m.streamErr = nil
			return m, m.zoteroItems.HelpText(items.ModeNormal)
		}
		m.loading = false
		m.streaming = true
		m.zoteroItems.AppendZoteroItems(msg.Items)
		return m, cmds.WaitForPageCmd(m.streamCh, m.streamErr, m.cache)

	case search.SearchMsg:
		m.loading = true
		m.currentView = ItemsView
		m.zoteroItems.ClearItems()
		return m, tea.Batch(m.spinner.Tick,
			cmds.StreamSearchCmd(m.zotero, msg.Query))

	case cmds.ChildrenLoadedMsg:
		if msg.Err != nil {
			return m, nil
		}
		m.zoteroItems.UpdateChildrenItems(msg.ParentKey, msg.Attachments, msg.Notes)
		return m, nil

	case attachpicker.AttachmentSelectedMsg:
		m.currentView = ItemsView
		cmds.OpenPDF(m.systemPDFOpener, msg.Key, msg.Filename)
		return m, nil

	case cmds.BibMsg:
		if msg.Err != nil {
			return m, nil
		}
		clipboard.Write(msg.Bib)
		return m, m.zoteroItems.HelpText(items.ModeClipboard)

	case cmds.ResetHelpMsg:
		m.zoteroItems.HelpText(items.ModeNormal)
		return m, nil
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	if m.loading || m.streaming {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch m.currentView {
	case InitialView:
		m.initial, cmd = m.initial.Update(msg)
	case CollectionsView:
		m.collections, cmd = m.collections.Update(msg)
	case ItemsView:
		m.zoteroItems, cmd = m.zoteroItems.Update(msg)
	case NoteEditorView:
		m.noteEditor, cmd = m.noteEditor.Update(msg)
	case NotePickerView:
		m.notepicker, cmd = m.notepicker.Update(msg)
	case SearchView:
		m.searchInput, cmd = m.searchInput.Update(msg)
	case AttachmentView:
		m.attachReader, cmd = m.attachReader.Update(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

var (
	headerStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true)
)

func (m rootModel) View() string {
	header := headerStyle.Width(m.width).Render("zui")

	var body string

	if m.loading {
		body = lipgloss.NewStyle().
			Height(m.height-3).
			Width(m.width).
			Align(lipgloss.Center, lipgloss.Center).
			Render(fmt.Sprintf("\n%s Loading Items", m.spinner.View()))
	} else {
		switch m.currentView {
		case InitialView:
			body = m.initial.View()
		case CollectionsView:
			body = m.collections.View()
		case ItemsView:
			body = m.zoteroItems.View()
		case NoteEditorView:
			body = m.noteEditor.View()
		case NotePickerView:
			body = m.notepicker.View()
		case SearchView:
			body = m.searchInput.View()
		case AttachmentView:
			body = m.attachReader.View()
		}
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
	)

	return content
}

func (m rootModel) isFiltering() bool {
	switch m.currentView {
	case CollectionsView:
		return m.collections.IsFiltering()
	case ItemsView:
		return m.zoteroItems.IsFiltering()
	}
	return false
}

func (m rootModel) isFilterApplied() bool {
	switch m.currentView {
	case CollectionsView:
		return m.collections.IsFilterApplied()
	case ItemsView:
		return m.zoteroItems.IsFilterApplied()
	}
	return false
}

func (m rootModel) externalEditor() string {
	if m.zotero.Config != nil {
		return m.zotero.Config.NoteEditor
	}
	return ""
}
