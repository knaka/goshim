package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getGoBinDir() (binDir string) {
	if len(os.Getenv("GOBIN")) > 0 {
		binDir = os.Getenv("GOBIN")
	} else if len(os.Getenv("GOPATH")) > 0 {
		binDir = filepath.Join(os.Getenv("GOPATH"), "bin")
	} else {
		home, err := os.UserHomeDir()
		panicOn(err)
		binDir = filepath.Join(home, "go", "bin")
	}
	fileinfo, err := os.Stat(binDir)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			panicOn(err)
		}
		err = os.MkdirAll(binDir, 0755)
		panicOn(err)
	}
	if !fileinfo.IsDir() {
		log.Fatalf("%s is not directory\n", binDir)
	}
	return
}

func calledAsGoshim(cmd string) bool {
	return strings.HasSuffix(cmd, "goshim")
}

func main() {
	waitForDebugger()

	err := createConfigFileIfNotExists(userConfigDir)
	panicOn(err)

	config, err := unmarshalConfigFile(filepath.Join(userConfigDir, "goshim.toml"))
	panicOn(err)

	if !calledAsGoshim(os.Args[0]) {
		_ = execCommandAndNotReturn(config, os.Args)
		panic("Failed to exec(2)")
	}

	pflag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("Usage: %v [options] subcmd", os.Args[0]))
		_, _ = fmt.Fprintln(os.Stderr)
		_, _ = fmt.Fprintln(os.Stderr, "Options:")
		pflag.PrintDefaults()
		_, _ = fmt.Fprintln(os.Stderr)
		_, _ = fmt.Fprintln(os.Stderr, "Subcmds:")
		_, _ = fmt.Fprintln(os.Stderr, "  install: Installs symlinks")
	}
	var debug bool
	pflag.BoolVarP(&debug, "debug", "", false, "debug")
	var shouldPutHelpAndExit bool
	pflag.BoolVarP(&shouldPutHelpAndExit, "help", "h", false, "display this help and exit")
	_ = pflag.CommandLine.Parse(os.Args[1:])
	if shouldPutHelpAndExit || len(pflag.Args()) == 0 {
		pflag.Usage()
		os.Exit(0)
	}

	cmdPath, err := os.Executable()
	panicOn(err)
	cmdPath, err = filepath.EvalSymlinks(cmdPath)
	panicOn(err)
	cmdPath = filepath.Clean(cmdPath)
	if debug {
		_, _ = fmt.Fprintln(os.Stderr, "cmdPath", cmdPath)
	}

	subCmd := pflag.Args()[0]
	switch subCmd {
	case "install":
		_ = updateSymlinks(cmdPath, config.Projects)
	case "rebuild":
		// todo
	}
}

func removeSysmlinks(cmdPath, binDir string) error {
	// If no result, cmdLinkPaths = nil and err = nil
	cmdLinkPaths, err := filepath.Glob(filepath.Join(binDir, "*"))
	if err != nil {
		return err
	}
	for _, cmdLinkPath := range cmdLinkPaths {
		fileInfo, err := os.Lstat(cmdLinkPath)
		if err != nil || (fileInfo.Mode()&os.ModeSymlink != os.ModeSymlink) {
			continue
		}
		linkTarget, err := filepath.EvalSymlinks(cmdLinkPath)
		if err != nil {
			continue
		}
		if cmdPath == linkTarget {
			err = os.Remove(cmdLinkPath)
			panicOn(err)
		}
	}
	return nil
}

func updateSymlinks(cmdPath string, projects []Project) error {
	binDir := getGoBinDir()
	_ = removeSysmlinks(cmdPath, binDir)
	for _, project := range projects {
		cmdSrcDirs, err := filepath.Glob(filepath.Join(project.Directory, "cmd", "*"))
		if err != nil {
			return err
		}
		for _, cmdSrcDir := range cmdSrcDirs {
			cmdBase := filepath.Base(cmdSrcDir)
			if cmdBase[0] == '_' {
				continue
			}
			linkPath := filepath.Join(binDir, cmdBase)
			_ = os.Remove(linkPath)
			err = os.Symlink(cmdPath, linkPath)
			panicOn(err)
		}
	}
	return nil
}
