package items

import (
	"fmt"
	"io"

	"github.com/camilo-zuluaga/zui/zotero"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#0F0F0F")).Background(lipgloss.Color("#F2F2F2"))
	normalTitle   = lipgloss.NewStyle()
	normalDesc    = lipgloss.NewStyle().Faint(true)
)

type item struct {
	title string
	desc  string

	ZoteroItem zotero.ZoteroItem
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type delegate struct{}

func NewDelegate() list.ItemDelegate {
	return delegate{}
}

func (d delegate) Height() int                               { return 2 }
func (d delegate) Spacing() int                              { return 1 }
func (d delegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d delegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	it, ok := listItem.(item)
	if !ok {
		return
	}

	var title, desc string
	if index == m.Index() {
		title = selectedStyle.Render(it.title)
		desc = normalDesc.Render(it.desc)
	} else {
		title = normalTitle.Render(it.title)
		desc = normalDesc.Render(it.desc)
	}

	fmt.Fprintf(w, "%s\n%s", title, desc)
}
