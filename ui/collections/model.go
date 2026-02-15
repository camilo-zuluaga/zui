package collections

import (
	"fmt"

	"github.com/camilo-zuluaga/zotero-tui/zotero"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C"))

type Model struct {
	height, width int
	list          list.Model
}

func New() Model {
	items := []list.Item{}

	l := list.New(items, NewDelegate(), 0, 0)
	l.SetStatusBarItemName("Collection", "Collections")
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.Styles.StatusBar = l.Styles.StatusBar.
		Align(lipgloss.Left).
		Margin(0).
		MarginBottom(1).
		Padding(0)

	return Model{list: l}
}

func (m *Model) SetZoteroCollections(zCollections []zotero.Collection) {
	items := make([]list.Item, 0, len(zCollections))
	for _, z := range zCollections {
		c := item{
			title:      z.Data.Name,
			desc:       fmt.Sprintf("%d items", z.Meta.NumItems),
			Collection: z,
		}
		items = append(items, c)
	}
	m.list.SetItems(items)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	content := lipgloss.NewStyle().
		Width(m.list.Width()).
		Align(lipgloss.Center).
		Render(m.list.View())

	help := helpStyle.Render("[s] Search items")

	return lipgloss.JoinVertical(lipgloss.Center, content, help)
}

func (m *Model) SetSize(width, height int) {
	m.list.SetSize(width-4, height-4)
}

func (m Model) SelectedCollection() *zotero.Collection {
	it, ok := m.list.SelectedItem().(item)
	if !ok {
		return nil
	}
	return &it.Collection
}
