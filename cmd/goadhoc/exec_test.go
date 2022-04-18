package main

import (
	"path/filepath"
	"testing"
)

// $ ls _test/sources/*.go | sort | xargs cat | shasum --algorithm 256
func TestSourcesHash(t *testing.T) {
	if sourcesHash(filepath.Join("_test", "sources", "*.go")) != "5791fc3" {
		t.Fatalf("not matched")
	}
}
