package main

import (
	"path/filepath"
	"testing"
)

func TestSourcesHash(t *testing.T) {
	result := getSourcesHash(filepath.Join("_test", "sources", "*.go"))
	// $ ls _test/sources/*.go | sort | xargs cat | shasum --algorithm 256
	if result != "24cd6e1" {
		t.Fatalf("not matched")
	}
}
