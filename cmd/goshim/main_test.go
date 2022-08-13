package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestRemoveSymlinks(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test")
	panicOn(err)
	defer func() { _ = os.RemoveAll(tempDir) }()
	removeSysmlinks("", tempDir)
}
