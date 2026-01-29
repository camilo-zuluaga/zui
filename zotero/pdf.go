package zotero

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var zoteroStoragePath = filepath.Join("Zotero", "storage")

type CommandRunner func(name string, path string) error

type SystemPDFOpener struct {
	SystemPath    string
	CommandOpener string
	RunCmd        CommandRunner
}

func NewSystemPDFOpener() *SystemPDFOpener {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Couldn't detect home directory")
	}
	zoteroPath := filepath.Join(homeDir, zoteroStoragePath)

	return &SystemPDFOpener{
		SystemPath:    zoteroPath,
		CommandOpener: detectPDFOpener(),
		RunCmd: func(name string, path string) error {
			return exec.Command(name, path).Run()
		},
	}
}

func (s *SystemPDFOpener) Open(collection, filename string) error {
	path := filepath.Join(s.SystemPath, collection, filename)
	return s.RunCmd(s.CommandOpener, path)
}

func detectPDFOpener() string {
	os := runtime.GOOS

	// TODO: implement windows one, although i dont know if this works in mac
	switch os {
	case "linux":
		return "xdg-open"
	case "darwin":
		return "open"
	default:
		return ""
	}
}
