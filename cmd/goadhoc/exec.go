package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"syscall"
)

func sourcesHash(wildcard string) string {
	files, err := filepath.Glob(wildcard)
	sort.Strings(files)
	panicOn(err)
	hasher := sha256.New()
	for _, file := range files {
		infile, _ := os.Open(file)
		_, _ = io.Copy(hasher, infile)
	}
	return hex.EncodeToString(hasher.Sum(nil))[0:7]
}

func execCommandAndNotReturn(config *appConfig, args []string) int {
	cmdBase := filepath.Base(args[0])
	/* cmdRealPath */ _, err := filepath.EvalSymlinks(args[0])
	var goSrcDir string
	var cmdDir string
	config.walkProjectCmds(func(project *Project, dir string) {
		base := filepath.Base(dir)
		if base == cmdBase {
			goSrcDir = project.Directory
			cmdDir = dir
		}
	})
	if goSrcDir == "" {
		panic("not found")
	}
	bindir := getBinDir()
	cacheDir := filepath.Join(bindir, ".goadhoc")
	err = os.MkdirAll(cacheDir, 0755)
	panicOn(err)
	srcHash := sourcesHash(filepath.Join(cmdDir, "*.go"))
	cacheBin := filepath.Join(cacheDir, fmt.Sprintf("%v.%v", cmdBase, srcHash))
	if _, err = os.Stat(cacheBin); err != nil {
		oldBins, err := filepath.Glob(filepath.Join(cacheDir, fmt.Sprintf("%v.*", cmdBase)))
		panicOn(err)
		for _, oldBin := range oldBins {
			_ = os.Remove(oldBin)
		}
		prevDir, _ := filepath.Abs(".")
		err = os.Chdir(goSrcDir)
		panicOn(err)
		pwd, _ := os.Getwd()
		log.Println("cp1:", pwd)
		b, err := exec.Command("go", "build", "-o", cacheBin, cmdDir).Output()
		if err != nil {
			log.Println(string(b))
			panic("compilation failed")
		}
		os.Chdir(prevDir)
	}
	args[0] = cacheBin
	err = syscall.Exec(cacheBin, args, os.Environ())
	return 1
}
