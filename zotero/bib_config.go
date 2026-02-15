package zotero

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type BibConfig struct {
	Format string
	Style  string
}

func LoadConfig() (*BibConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Couldn't detect home directory")
	}
	path := filepath.Join(homeDir, ".config", "zui", "config.toml")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("Config file does not exist.")
	} else if err != nil {
		return nil, err
	}

	var conf BibConfig
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
