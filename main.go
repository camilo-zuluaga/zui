package main

import (
	"fmt"
	"os"

	"github.com/99designs/keyring"
	"github.com/camilo-zuluaga/zotero-tui/cache"
	"github.com/camilo-zuluaga/zotero-tui/sync"
	"github.com/camilo-zuluaga/zotero-tui/ui"
	"github.com/camilo-zuluaga/zotero-tui/zotero"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	db, err := cache.Init()
	if err != nil {
		fmt.Println("failed to init cache:", err)
		os.Exit(1)
	}
	defer db.Close()

	kr, err := keyring.Open(keyring.Config{
		ServiceName: "zotero-tui",
	})

	var model tea.Model
	if err == nil {
		apiItem, apiErr := kr.Get("api-key")
		idItem, idErr := kr.Get("user-id")
		if apiErr == nil && idErr == nil {
			zclient := zotero.NewZoteroClient("https://api.zotero.org", string(idItem.Data), string(apiItem.Data))
			ss := sync.New(db, zclient)
			model = ui.NewRootModel(zclient, db, ss)
		}
	}

	if model == nil {
		model = ui.NewInitialRootModel(db)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
