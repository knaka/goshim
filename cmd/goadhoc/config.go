package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//go:embed goadhoc-default.toml
var defaultConfigToml []byte

const homeVariable = "$HOME"

type Project struct {
	Directory string `toml:"directory"`
}

type appConfig struct {
	Foo      string    `toml:"foo"`
	Projects []Project `toml:"projects"`
}

func (config *appConfig) walkProjectCmds(walker func(*Project, string)) {
	for _, project := range config.Projects {
		dirs, err := filepath.Glob(filepath.Join(project.Directory, "cmd", "*"))
		panicOn(err)
		for _, dir := range dirs {
			walker(&project, dir)
		}
	}
}

func createConfigFileIfNotExists(confDir string) (err error) {
	confPath := filepath.Join(confDir, "goadhoc.toml")
	if _, err := os.Stat(confPath); err == nil {
		return err
	}
	err = os.MkdirAll(confDir, 0755)
	if err != nil {
		return err
	}
	outfile, err := os.OpenFile(confPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	ingsReader := bytes.NewReader(defaultConfigToml)
	_, err = io.Copy(outfile, ingsReader)
	if err != nil {
		return err
	}
	return nil
}

func unmarshalConfigFile(confPath string) (*appConfig, error) {
	data, _ := ioutil.ReadFile(confPath)
	config := &appConfig{}
	err := toml.Unmarshal([]byte(data), config)
	homeDir, err := os.UserHomeDir()
	panicOn(err)
	homeDir = filepath.Clean(homeDir)
	homeVariableWithSeparator := fmt.Sprintf("%v%c", homeVariable, filepath.Separator)
	for i, project := range config.Projects {
		directory := strings.Replace(project.Directory, homeVariable, homeDir, 1)
		if strings.Index(project.Directory, homeVariableWithSeparator) == 0 {
			directory = strings.Replace(project.Directory, homeVariable, homeDir, 1)
		}
		directory, err = filepath.EvalSymlinks(directory)
		panicOn(err)
		config.Projects[i].Directory = directory
	}
	return config, err
}
