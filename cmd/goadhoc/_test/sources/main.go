package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getBinDir() string {
	var binDir string
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
	return binDir
}

func isCalledAsGoadhoc(cmd string) bool {
	return strings.HasSuffix(cmd, "goadhoc")
}

func main() {
	waitForDebugger()

	homeDir, err := os.UserHomeDir()
	panicOn(err)
	// todo: 配慮 for the the other platforms
	confDir := filepath.Join(homeDir, ".config")
	err = createConfigFileIfNotExists(confDir)
	panicOn(err)

	config, err := unmarshalConfigFile(filepath.Join(confDir, "goadhoc.toml"))
	panicOn(err)

	if !isCalledAsGoadhoc(os.Args[0]) {
		os.Exit(runCommand(config, os.Args))
	}

	pflag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [options]:\n", os.Args[0])
		pflag.PrintDefaults()
	}
	var debug bool
	pflag.BoolVarP(&debug, "debug", "", false, "debug")
	var shouldPutHelpAndExit bool
	pflag.BoolVarP(&shouldPutHelpAndExit, "help", "h", false, "display this help and exit")
	_ = pflag.CommandLine.Parse(os.Args[1:])
	if shouldPutHelpAndExit {
		pflag.Usage()
		os.Exit(0)
	}
	if len(pflag.Args()) == 0 {
		log.Fatal("too few argument")
	}
	cmdPath, err := filepath.Abs(os.Args[0])
	panicOn(err)
	cmdPath = filepath.Clean(cmdPath)
	if debug {
		_, _ = fmt.Fprintln(os.Stderr, "cmdPath", cmdPath)
	}

	for _, project := range config.Projects {
		fmt.Println("Project directory", project.Directory)
	}

	if debug {
		log.Println("0", os.Args[0])
	}

	cmd := pflag.Args()[0]
	switch cmd {
	case "install":
		_ = updateSymlinks(cmdPath, config.Projects)
	case "update":
		_ = updateSymlinks(cmdPath, config.Projects)
		// todo:
		// compileAll(config.Projects)
	}
}

func removeSysmlinks(cmdPath, bindir string) error {
	// If no result, cmds = nil and err = nil
	cmds, err := filepath.Glob(filepath.Join(bindir, "*"))
	if err != nil {
		return err
	}
	for _, cmd := range cmds {
		target, err := filepath.EvalSymlinks(cmd)
		panicOn(err)
		// todo:
		if cmdPath == cmd {
			fmt.Fprintln(os.Stderr, "removing", target)
		}
	}
	return nil
}

func updateSymlinks(cmdPath string, projects []Project) error {
	bindir := getBinDir()
	_ = removeSysmlinks(cmdPath, bindir)
	for _, project := range projects {
		cmdDirs, err := filepath.Glob(filepath.Join(project.Directory, "cmd", "*"))
		if err != nil {
			return err
		}
		for _, cmdDir := range cmdDirs {
			base := filepath.Base(cmdDir)
			if base[0] == '_' {
				continue
			}
			linkPath := filepath.Join(bindir, base)
			err = os.Symlink(cmdPath, linkPath)
			panicOn(err)
		}
	}
	return nil
}
