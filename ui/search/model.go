package search

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C"))

type SearchMsg struct {
	Query string
}

type Model struct {
	textInput textinput.Model
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Attention is all you need"
	ti.Focus()
	ti.Width = 100

	return Model{
		textInput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			return m, func() tea.Msg {
				return SearchMsg{
					Query: m.textInput.Value(),
				}
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"Input the [Title | Creator | Year] to search:\n\n%s\n\n%s",
		m.textInput.View(),
		helpStyle.Render("(esc to quit)"),
	) + "\n"
}
