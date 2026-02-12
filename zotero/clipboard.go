package zotero

import (
	"fmt"
	"os"

	"golang.design/x/clipboard"
)

func InitClipboard() {
	if err := clipboard.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Clipboard init error: %v\n", err)
		os.Exit(1)
	}
}

func WriteClipboard(text string) {
	clipboard.Write(clipboard.FmtText, []byte(text))
}
