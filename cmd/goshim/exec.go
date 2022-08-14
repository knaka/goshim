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

func getSourcesHash(wildcard string) string {
	paths, err := filepath.Glob(wildcard)
	panicOn(err)
	sort.Strings(paths)
	hashGenerator := sha256.New()
	for _, path := range paths {
		func() {
			infile, err := os.Open(path)
			panicOn(err)
			defer func() { _ = infile.Close() }()
			_, err = io.Copy(hashGenerator, infile)
			panicOn(err)
		}()
	}
	return hex.EncodeToString(hashGenerator.Sum(nil))[0:7]
}

func execCommandAndNotReturn(config *appConfig, args []string) int {
	cmdBase := filepath.Base(args[0])
	var goProjectDir string
	var srcDir string
	config.walkProjectCmds(
		func(project *Project, srcDirCandidate string) (finished bool) {
			srcDirBase := filepath.Base(srcDirCandidate)
			if cmdBase == srcDirBase {
				goProjectDir = project.Directory
				srcDir = srcDirCandidate
				return true
			}
			return false
		},
	)
	if goProjectDir == "" || srcDir == "" {
		panic("Source dir not found")
	}
	binDir := getGoBinDir()
	cacheDir := filepath.Join(binDir, ".goshim")
	err := os.MkdirAll(cacheDir, 0755)
	panicOn(err)
	hash := getSourcesHash(filepath.Join(srcDir, "*.go"))
	cacheBinPath := filepath.Join(cacheDir, fmt.Sprintf("%v.%v", cmdBase, hash))
	if _, err = os.Stat(cacheBinPath); err != nil {
		oldBinPaths, err := filepath.Glob(filepath.Join(cacheDir, fmt.Sprintf("%v.*", cmdBase)))
		panicOn(err)
		for _, oldBinPath := range oldBinPaths {
			_ = os.Remove(oldBinPath)
		}
		savedDir, _ := filepath.Abs(".")
		err = os.Chdir(goProjectDir)
		panicOn(err)
		output, err := exec.Command("go", "build", "-o", cacheBinPath, srcDir).Output()
		if err != nil {
			log.Println(string(output))
			panic("Compilation failed")
		}
		_ = os.Chdir(savedDir)
	}
	args[0] = cacheBinPath
	err = syscall.Exec(cacheBinPath, args, os.Environ())
	return 1
}
