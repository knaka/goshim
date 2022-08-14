//go:build !windows

package main

import (
	"os"
	"path/filepath"
)

var userConfigDir string

func init() {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		userConfigDir = os.Getenv("XDG_CONFIG_HOME")
	} else {
		homeDir, err := os.UserHomeDir()
		panicOn(err)
		userConfigDir = filepath.Join(homeDir, ".config")
	}
}
