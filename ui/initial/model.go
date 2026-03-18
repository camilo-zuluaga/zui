package initial

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C"))
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#E0E0E0"))
	urlStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#7EB8DA")).Underline(true)
	stepStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0"))
	bulletStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7EB8DA"))
)

type CredentialsMsg struct {
	APIKey string
	UserID string
}

type Model struct {
	textInputAPI textinput.Model
	textInputID  textinput.Model
	step         int // 0 = API key, 1 = user ID
}

func InitialModel() Model {
	tiAPI := textinput.New()
	tiAPI.Placeholder = "Copy your Zotero API key here"
	tiAPI.Focus()
	tiAPI.Width = 100

	tiID := textinput.New()
	tiID.Placeholder = "Copy your Zotero user ID here"
	tiID.Blur()
	tiID.Width = 100
	return Model{
		textInputAPI: tiAPI,
		textInputID:  tiID,
		step:         0,
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
			if m.step == 0 {
				m.step = 1
				m.textInputAPI.Blur()
				m.textInputID.Focus()
				return m, textinput.Blink
			}
			return m, func() tea.Msg {
				return CredentialsMsg{
					APIKey: m.textInputAPI.Value(),
					UserID: m.textInputID.Value(),
				}
			}
		}
	}

	if m.step == 0 {
		m.textInputAPI, cmd = m.textInputAPI.Update(msg)
	} else {
		m.textInputID, cmd = m.textInputID.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	b := bulletStyle.Render("•")
	apiKeyHelp := fmt.Sprintf(
		"%s\n\n  %s %s\n  %s %s\n  %s %s\n  %s\n  %s %s\n\n %s\n\n  %s %s\n %s\n\n %s\n %s\n",
		titleStyle.Render("  How to get your Zotero API key:"),
		b, stepStyle.Render("Go to "+urlStyle.Render("https://www.zotero.org/settings/keys/new")),
		b, stepStyle.Render("Set a name for your key."),
		b, stepStyle.Render("The permissions you need to check:"),
		stepStyle.Render("	1. Allow library access\n      2. Allow notes access\n      3. Allow write access\n      4. Read/Write on Group permissions"),
		b, stepStyle.Render("Save the key and copy the generated key."),
		titleStyle.Render(" How to get your Zotero user ID:"),
		b, stepStyle.Render("Go to "+urlStyle.Render("https://www.zotero.org/settings/security#applications")),
		stepStyle.Render("Your user ID is shown on that section of the page."),
		bulletStyle.Render("The text shows as:"),
		stepStyle.Render("\"User ID: Your user ID for use in API calls is 1xxxxxxx\""),
	)

	var input string
	if m.step == 0 {
		input = fmt.Sprintf("  API Key:\n  %s", m.textInputAPI.View())
	} else {
		input = fmt.Sprintf("  API Key: %s\n\n  User ID:\n  %s", m.textInputAPI.Value(), m.textInputID.View())
	}

	return fmt.Sprintf(
		"Hello! in order to use zui you will need to set your Zotero API key and your user ID:\n\n%s\n\n%s\n",
		apiKeyHelp,
		input,
	) + "\n"
}
