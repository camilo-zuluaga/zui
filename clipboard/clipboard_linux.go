//go:build linux || freebsd || openbsd || netbsd || dragonfly

package clipboard

import (
	"os"
	"os/exec"
	"strings"
)

func Write(text string) {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		cmd := exec.Command("wl-copy")
		cmd.Stdin = strings.NewReader(text)
		cmd.Run()
		return
	}

	cmd := exec.Command("xclip", "-selection", "clipboard")
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err == nil {
		return
	}

	cmd = exec.Command("xsel", "--clipboard", "--input")
	cmd.Stdin = strings.NewReader(text)
	cmd.Run()
}
