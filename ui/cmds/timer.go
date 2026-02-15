package cmds

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type ResetHelpMsg struct{}

func ResetHelpCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		return ResetHelpMsg{}
	}
}
