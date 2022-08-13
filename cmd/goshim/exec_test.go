package main

import (
	"path/filepath"
	"testing"
)

// $ ls _test/sources/*.go | sort | xargs cat | shasum --algorithm 256
func TestSourcesHash(t *testing.T) {
	result := sourcesHash(filepath.Join("_test", "sources", "*.go"))
	if result != "24cd6e1" {
		t.Fatalf("not matched")
	}
}
