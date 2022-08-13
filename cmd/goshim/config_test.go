package main

import (
	"bytes"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func funcName(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func fileHash(path string) []byte {
	hasher := sha256.New()
	infile, _ := os.Open(path)
	_, _ = io.Copy(hasher, infile)
	return hasher.Sum(nil)
}

func Test_createConfigFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test")
	panicOn(err)
	defer func() { _ = os.RemoveAll(tempDir) }()

	basicTestDir := filepath.Join(tempDir, "basic_test")

	type args struct {
		confDir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		verify  func(t *testing.T)
	}{
		{
			"Basic",
			args{
				basicTestDir,
			},
			false,
			func(t *testing.T) {
				path := filepath.Join(basicTestDir, "goshim.toml")
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("File does not exist")
				}
				hashGotten := fileHash(path)
				hashExpected := fileHash(filepath.Join("_test", "config_basic_expected", "goshim.toml"))
				if bytes.Compare(hashGotten, hashExpected) != 0 {
					t.Fatalf("Files do not match")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createConfigFileIfNotExists(tt.args.confDir); (err != nil) != tt.wantErr {
				t.Fatalf("%v() error = %v, wantErr %v", funcName(createConfigFileIfNotExists), err, tt.wantErr)
			}
			tt.verify(t)
		})
	}
}
