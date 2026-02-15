package notepicker

import (
	"github.com/camilo-zuluaga/zotero-tui/ui/items"
	"github.com/camilo-zuluaga/zotero-tui/zotero"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)

type NoteSelectedMsg struct {
	ItemKey   string
	ParentKey string
	Content   string
	New       bool
}

type Model struct {
	height, width int
	list          list.Model

	parentKey string
}

func New(parentKey string) Model {
	notes := []list.Item{}

	l := list.New(notes, NewDelegate(), 0, 0)
	l.Title = titleStyle.Render("Select the note to edit.")
	l.SetStatusBarItemName("Note", "Notes")
	l.SetShowHelp(false)

	l.Styles.Title = l.Styles.Title.
		Align(lipgloss.Left).
		Margin(0).
		MarginBottom(1).
		Padding(0).
		Background(lipgloss.NoColor{})

	l.Styles.TitleBar = l.Styles.TitleBar.
		Align(lipgloss.Left).
		Margin(0).
		Padding(0).
		Background(lipgloss.NoColor{})

	l.Styles.StatusBar = l.Styles.StatusBar.
		Align(lipgloss.Left).
		Margin(0).
		MarginBottom(1).
		Padding(0)

	return Model{list: l, parentKey: parentKey}
}

func (m *Model) SetZoteroNotes(zNotes []zotero.ZoteroNote) {
	notes := make([]list.Item, 0, len(zNotes)+1)
	notes = append(notes, item{
		title: lipgloss.NewStyle().Faint(true).Render("󰎝 Create a note"),
	})
	for _, z := range zNotes {
		c := item{
			title: items.StripHTML(z.Note) + "...",
			note:  z,
		}
		notes = append(notes, c)
	}
	m.list.SetItems(notes)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if selected := m.selectedNote(); selected != nil {
				if selected.Note == "" {
					return m, func() tea.Msg {
						return NoteSelectedMsg{ParentKey: m.parentKey, New: true}
					}
				}

				return m, func() tea.Msg {
					return NoteSelectedMsg{
						ItemKey: selected.Key,
						Content: selected.Note,
						New:     false,
					}
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	content := lipgloss.NewStyle().
		PaddingLeft(1).
		Width(m.list.Width()).
		Render(m.list.View())

	return content
}

func (m *Model) SetSize(width, height int) {
	m.list.SetSize(width-4, height-4)
}

func (m Model) selectedNote() *zotero.ZoteroNote {
	it, ok := m.list.SelectedItem().(item)
	if !ok {
		return nil
	}
	return &it.note
}
