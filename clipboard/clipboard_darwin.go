//go:build darwin

package clipboard

import (
	"os/exec"
	"strings"
)

func Write(text string) {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	cmd.Run()
}
