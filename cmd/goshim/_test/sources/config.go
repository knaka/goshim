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

//go:embed goshim-default.toml
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
	confPath := filepath.Join(confDir, "goshim.toml")
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
		if strings.Index(project.Directory, homeVariableWithSeparator) != 0 {
			continue
		}
		config.Projects[i].Directory = strings.Replace(project.Directory, homeVariable, homeDir, 1)
	}
	return config, err
}

func (config appConfig) marshalConfigFile(confPath string) error {
	homeDir, err := os.UserHomeDir()
	panicOn(err)
	homeDir = filepath.Clean(homeDir)
	homeDirWithSeparator := fmt.Sprintf("%v%c", homeDir, filepath.Separator)
	for i, project := range config.Projects {
		if strings.Index(project.Directory, homeDirWithSeparator) != 0 {
			continue
		}
		config.Projects[i].Directory = strings.Replace(project.Directory, homeDir, homeVariable, 1)
	}
	s, err := toml.Marshal(config)
	if err != nil {
		return err
	}
	outfile, err := os.OpenFile(confPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	_, err = outfile.Write(s)
	if err != nil {
		return err
	}
	return nil
}
