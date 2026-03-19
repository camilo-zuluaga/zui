package attachpicker

import (
	"github.com/camilo-zuluaga/zui/zotero"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)

type AttachmentSelectedMsg struct {
	Key      string
	Filename string
}

type Model struct {
	itemTitle     string
	height, width int
	list          list.Model
}

func New(itemTitle string) Model {
	notes := []list.Item{}

	l := list.New(notes, NewDelegate(), 0, 0)
	l.SetShowTitle(false)
	l.SetStatusBarItemName("PDF", "PDFs")
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

	return Model{itemTitle: itemTitle, list: l}
}

func (m *Model) SetZoteroAttachments(za []zotero.ZoteroAttachment) {
	attachments := make([]list.Item, 0, len(za))
	for _, z := range za {
		attachment := item{
			title:      z.Title,
			attachment: z,
		}
		attachments = append(attachments, attachment)
	}
	m.list.SetItems(attachments)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if selected := m.selectedAttachment(); selected != nil {
				return m, func() tea.Msg {
					return AttachmentSelectedMsg{
						Key:      selected.Key,
						Filename: selected.Filename,
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
	title := titleStyle.Render(m.itemTitle)
	content := lipgloss.NewStyle().
		PaddingLeft(1).
		Width(m.list.Width()).
		Render(m.list.View())

	return lipgloss.JoinVertical(lipgloss.Left, title, content)
}

func (m *Model) SetSize(width, height int) {
	m.list.SetSize(width-4, height-4)
}

func (m Model) selectedAttachment() *zotero.ZoteroAttachment {
	it, ok := m.list.SelectedItem().(item)
	if !ok {
		return nil
	}
	return &it.attachment
}
