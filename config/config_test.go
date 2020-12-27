package config

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type projectConfig struct {
	projectpath        string
	projectname        string
	selectProjectIsNew bool
}

func (pc *projectConfig) GetProjectName() string {
	return pc.projectname
}

func (pc *projectConfig) GetProjectPath() string {
	return pc.projectpath
}

func (pc *projectConfig) SetProjectName(v string) {
	pc.projectname = v
}

func (pc *projectConfig) SetProjectPath(v string) {
	pc.projectpath = v
}

type testFileSystem struct {
	fileExists      func(string) (bool, error)
	dirExists       func(string) (bool, error)
	saveConfigFile  func(interface{}, string) error
	readConfigFile  func(string, interface{}) error
	goToProjectPath func(string) error
}

func (tfs *testFileSystem) FileExists(file string) (bool, error) {
	return tfs.fileExists(file)
}

func (tfs *testFileSystem) DirExists(path string) (bool, error) {
	return tfs.dirExists(path)
}

func (tfs *testFileSystem) SaveConfigFile(data interface{}, fileName string) error {
	return tfs.saveConfigFile(data, fileName)
}

func (tfs *testFileSystem) ReadConfigFile(fileName string, configuration interface{}) error {
	return tfs.readConfigFile(fileName, configuration)
}

func (tfs *testFileSystem) GoToProjectPath(projectpath string) error {
	return tfs.goToProjectPath(projectpath)
}

func TestLookupProjectConfigCase1(t *testing.T) {

	fs := testFileSystem{
		readConfigFile: func(filename string, configuration interface{}) error {
			return nil
		},
	}

	cfg := Config{}
	cfg.SetFileSystem(&fs)
	assert.Equal(t, nil, cfg.lookupProjectConfig())
}

func TestLookupProjectConfigCase2(t *testing.T) {
	fs := testFileSystem{
		readConfigFile: func(filename string, configuration interface{}) error {
			return errors.New("Cannot read file")
		},
	}

	cfg := Config{}
	cfg.SetFileSystem(&fs)
	err := cfg.lookupProjectConfig()
	assert.EqualError(t, err, "Cannot read file")
}

func TestLookupUserConfigCase1(t *testing.T) {
	fs := testFileSystem{
		readConfigFile: func(filename string, configuration interface{}) error {
			data := `{
	"projects": [
		{
			"path": "/path1/",
			"name": "project1",
			"main_container": "container_name1"
		},
		{
			"path": "",
			"name": "project3"
		},
		{
			"path": "",
			"name": ""
		},
		{
			"path": "/path2/",
			"name": "project2",
			"main_container": "container_name2"
		}
	],
	"settings": [
		{
			"path": "/path1/",
			"name": "project1",
			"main_container": "container_name1"
		},
		{
			"path": "",
			"name": "project3"
		}
	]
}`
			return json.Unmarshal([]byte(data), &configuration)
		},
	}

	cfg := Config{}
	cfg.Init()
	cfg.SetFileSystem(&fs)
	assert.Equal(t, nil, cfg.lookupUserConfig())

	assert.Equal(t, 4, len(cfg.globalConfig.Projects))
}

func TestLookupUserConfigCase2(t *testing.T) {
	fs := testFileSystem{
		saveConfigFile: func(data interface{}, fileName string) error {
			return errors.New("Check: if saveConfigFile has been called")
		},
		readConfigFile: func(filename string, configuration interface{}) error {
			return errors.New("error: no such file or directory")
		},
	}

	cfg := Config{}
	cfg.SetFileSystem(&fs)
	err := cfg.lookupUserConfig()
	assert.EqualError(t, err, "Check: if saveConfigFile has been called")
}

func TestLookupUserConfigCase3(t *testing.T) {
	fs := testFileSystem{
		readConfigFile: func(filename string, configuration interface{}) error {
			return errors.New("some different error")
		},
	}

	cfg := Config{}
	cfg.SetFileSystem(&fs)
	err := cfg.lookupUserConfig()
	assert.EqualError(t, err, "some different error")
}

func TestLoadConfigCase1(t *testing.T) {
	fs := testFileSystem{
		readConfigFile: func(filename string, configuration interface{}) error {
			return errors.New("user config error")
		},
		saveConfigFile: func(data interface{}, fileName string) error {
			return nil
		},
	}

	cfg := Config{}
	cfg.SetFileSystem(&fs)
	err := cfg.LoadConfig(true)
	assert.EqualError(t, err, "user config error")
}

func TestLoadConfigCase2(t *testing.T) {
	cnt := 0
	fs := testFileSystem{
		readConfigFile: func(filename string, configuration interface{}) error {

			if cnt == 0 {
				cnt++
				return nil
			}
			return errors.New("lookupProjectConfig error")
		},
		saveConfigFile: func(data interface{}, fileName string) error {
			return nil
		},
	}

	cfg := Config{}
	cfg.SetFileSystem(&fs)
	err := cfg.LoadConfig(true)
	assert.EqualError(t, err, "lookupProjectConfig error")
}

func TestLoadConfigCase3(t *testing.T) {
	cnt := 0
	fs := testFileSystem{
		readConfigFile: func(filename string, configuration interface{}) error {

			if cnt == 0 {
				cnt++
				return nil
			}
			return errors.New("error: no such file or directory path")
		},
		saveConfigFile: func(data interface{}, fileName string) error {
			return nil
		},
	}

	cfg := Config{}
	cfg.SetFileSystem(&fs)
	assert.Equal(t, nil, cfg.LoadConfig(true))
}

func TestFindProjectPathInJSONCase1(t *testing.T) {
	fs := testFileSystem{
		dirExists: func(path string) (bool, error) {
			if path == "path2" {
				return true, nil
			}
			return false, nil
		},
	}
	cfg := Config{}
	cfg.SetFileSystem(&fs)
	cfg.globalConfig = &GlobalConfig{
		Projects: []GlobalProjectConfig{
			GlobalProjectConfig{
				Path: "path1",
				Name: "name1",
			},
			GlobalProjectConfig{
				Path: "path2",
				Name: "name2",
			},
			GlobalProjectConfig{
				Path: "path3",
				Name: "name3",
			},
		},
	}

	pc := projectConfig{
		projectname: "name2",
	}

	cfg.FindProjectPathInJSON(&pc)

	assert.Equal(t, "path2", pc.GetProjectPath())
}

func TestFindProjectPathInJSONCase2(t *testing.T) {
	fs := testFileSystem{
		dirExists: func(path string) (bool, error) {
			if path == "path2" {
				return true, errors.New("Unknown error")
			}
			return false, nil
		},
	}
	cfg := Config{}
	cfg.SetFileSystem(&fs)
	cfg.globalConfig = &GlobalConfig{
		Projects: []GlobalProjectConfig{
			GlobalProjectConfig{
				Path: "path1",
				Name: "name1",
			},
			GlobalProjectConfig{
				Path: "path2",
				Name: "name2",
			},
			GlobalProjectConfig{
				Path: "path3",
				Name: "name3",
			},
		},
	}

	pc := projectConfig{
		projectname: "name2",
	}

	cfg.FindProjectPathInJSON(&pc)

	assert.Equal(t, "", pc.GetProjectPath())
}
