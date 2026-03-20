package zotero

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.design/x/clipboard"
)

func InitClipboard() {
	if err := clipboard.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Clipboard init error: %v\n", err)
		os.Exit(1)
	}
}

func WriteClipboard(text string) {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		cmd := exec.Command("wl-copy")
		cmd.Stdin = strings.NewReader(text)
		cmd.Run()
		return
	}
	clipboard.Write(clipboard.FmtText, []byte(text))
}
