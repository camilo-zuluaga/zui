//go:build windows

package clipboard

import (
	"os/exec"
	"strings"
)

func Write(text string) {
	cmd := exec.Command("clip.exe")
	cmd.Stdin = strings.NewReader(text)
	cmd.Run()
}
