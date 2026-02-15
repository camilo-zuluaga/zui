package noteeditor

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C"))

type CancelNoteMsg struct{}

type SavedNoteMsg struct {
	ParentKey string
	Key       string
	Content   string
	New       bool
}

type Model struct {
	width, height int
	textarea      textarea.Model

	itemKey   string
	parentKey string
	new       bool
}

func InitialModel(parentKey, itemKey, content string, isNew bool) Model {
	ti := textarea.New()
	ti.Focus()

	if !isNew {
		ti.SetValue(content)
	}

	return Model{
		textarea:  ti,
		itemKey:   itemKey,
		parentKey: parentKey,
		new:       isNew,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case "ctrl+c":
			return m, func() tea.Msg { return CancelNoteMsg{} }
		case "ctrl+s":
			return m, func() tea.Msg {
				return SavedNoteMsg{
					ParentKey: m.parentKey,
					Key:       m.itemKey,
					Content:   m.textarea.Value(),
					New:       m.new,
				}
			}
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return fmt.Sprintf(
		"Write your note.\n\n%s\n\n%s",
		m.textarea.View(),
		helpStyle.Render("[ctrl + s] Save • [ctrl + c] Cancel"),
	) + "\n\n"
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	m.textarea.SetWidth(width - 2)
	m.textarea.SetHeight(height - 8)
}
