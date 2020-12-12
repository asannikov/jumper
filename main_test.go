package main

import (
	"encoding/json"
	"errors"
	"mgt/config"
	"mgt/dialog"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TestProjectFileContent = `{
	"path": "/path/to/project/",
	"name": "project name",
	"main_container": ""
}`

const TestUserFileContent = `{
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

type testFileSystem struct {
	fileExists       func(string) (bool, error)
	dirExists        func(string) (bool, error)
	saveConfigFile   func(interface{}, string) error
	readConfigFile   func(string, interface{}) error
	goToProjectPath  func(string) error
	getUserDirectory func() (string, error)
	getWd            func() (string, error)
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

func (tfs *testFileSystem) GetUserDirectory() (string, error) {
	return tfs.getUserDirectory()
}
func (tfs *testFileSystem) GetWd() (string, error) {
	return tfs.getWd()
}

func TestDefinePaths(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	tfs := &testFileSystem{
		getUserDirectory: func() (string, error) {
			return "/user/path/", nil
		},
	}

	err := definePaths(cfg, tfs)

	assert.Equal(t, nil, err)
	assert.Equal(t, "/user/path/"+string(os.PathSeparator)+".mgt.json", cfg.UserFile)

	tfs = &testFileSystem{
		getUserDirectory: func() (string, error) {
			return "", errors.New("Cannot read directory")
		},
	}

	err = definePaths(cfg, tfs)
	assert.EqualError(t, err, "Cannot read directory")
}

func TestSeekPath_DefinePaths(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	tfs := &testFileSystem{
		getUserDirectory: func() (string, error) {
			return "", errors.New("Cannot read directory")
		},
	}

	dialog := &dialog.Dialog{}

	err := seekPath(cfg, dialog, tfs)
	assert.EqualError(t, err, "Cannot read directory")
}

func TestSeekPath_LoadConfig(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	tfs := &testFileSystem{
		getUserDirectory: func() (string, error) {
			return "/user/path/", nil
		},
		readConfigFile: func(filename string, configuration interface{}) error {
			return errors.New("Cannot read config file")
		},
	}

	cfg.SetFileSystem(tfs)

	dialog := &dialog.Dialog{}

	err := seekPath(cfg, dialog, tfs)
	assert.EqualError(t, err, "Cannot read config file")
}

func TestSeekPath_GetWd(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	tfs := &testFileSystem{
		getUserDirectory: func() (string, error) {
			return "/user/path/", nil
		},
		readConfigFile: func(filename string, configuration interface{}) error {
			return nil
		},
		getWd: func() (string, error) {
			return "", errors.New("Cannot get current directory")
		},
	}

	cfg.SetFileSystem(tfs)

	dialog := &dialog.Dialog{}

	err := seekPath(cfg, dialog, tfs)
	assert.EqualError(t, err, "Cannot get current directory")
}

func TestSeekPath_GetProjectNameList(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	tfs := &testFileSystem{
		getUserDirectory: func() (string, error) {
			return "/user/path/", nil
		},
		saveConfigFile: func(data interface{}, fileName string) error {
			return nil
		},
		getWd: func() (string, error) {
			return "/current/path/", errors.New("stop execution")
		},
		readConfigFile: func(filename string, configuration interface{}) error {
			if filename == cfg.UserFile {
				return json.Unmarshal([]byte(TestUserFileContent), &configuration)
			} else if filename == cfg.ProjectFile {
				return json.Unmarshal([]byte(TestProjectFileContent), &configuration)
			}
			return errors.New("Wrong test file")
		},
	}

	cfg.SetFileSystem(tfs)
	dialog := &dialog.Dialog{}

	err := seekPath(cfg, dialog, tfs)
	assert.EqualError(t, err, "stop execution")

	pl := cfg.GetProjectNameList()
	assert.Equal(t, 2, len(pl))
}

func TestSeekPath_runDialogCase1(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	tfs := &testFileSystem{
		getUserDirectory: func() (string, error) {
			return "/user/path/", nil
		},
		saveConfigFile: func(data interface{}, fileName string) error {
			return nil
		},
		getWd: func() (string, error) {
			return "/current/path/", nil
		},
		readConfigFile: func(filename string, configuration interface{}) error {
			if filename == cfg.UserFile {
				return json.Unmarshal([]byte(TestUserFileContent), &configuration)
			} else if filename == cfg.ProjectFile {
				return errors.New("Error: no such file or directory")
			}
			return errors.New("Wrong test file")
		},
	}

	cfg.SetFileSystem(tfs)
	dialog := &dialog.Dialog{
		SelectProject: func(projects []string) (int, string, error) {
			return 0, "", errors.New("Error SelectProject dialog")
		},
	}

	err := seekPath(cfg, dialog, tfs)
	assert.EqualError(t, err, "Error SelectProject dialog")
}

func TestSeekPath_runDialogCase2(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	tfs := &testFileSystem{
		getUserDirectory: func() (string, error) {
			return "/user/path/", nil
		},
		saveConfigFile: func(data interface{}, fileName string) error {
			return nil
		},
		getWd: func() (string, error) {
			return "/current/path/", nil
		},
		readConfigFile: func(filename string, configuration interface{}) error {
			if filename == cfg.UserFile {
				return json.Unmarshal([]byte(TestUserFileContent), &configuration)
			} else if filename == cfg.ProjectFile {
				return errors.New("Error: no such file or directory")
			}
			return errors.New("Wrong test file")
		},
		goToProjectPath: func(path string) error {
			return nil
		},
	}

	cfg.SetFileSystem(tfs)
	dialog := &dialog.Dialog{
		SelectProject: func(projects []string) (int, string, error) {
			return -1, "New Project Name", nil
		},
		AddProjectName: func() (string, error) {
			return "", errors.New("Should not be called as the project name has been typed")
		},
		AddProjectPath: func(path string) (string, error) {
			if path == "/current/path/" {
				return path, errors.New("Is ok")
			}
			return "", errors.New("Current path was not set")
		},
	}

	err := seekPath(cfg, dialog, tfs)
	assert.EqualError(t, err, "Is ok")
}

func TestSeekPath_runDialogCase3(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	tfs := &testFileSystem{
		getUserDirectory: func() (string, error) {
			return "/user/path/", nil
		},
		saveConfigFile: func(data interface{}, fileName string) error {
			return nil
		},
		getWd: func() (string, error) {
			return "/current/path/", nil
		},
		readConfigFile: func(filename string, configuration interface{}) error {
			if filename == cfg.UserFile {
				return json.Unmarshal([]byte(TestUserFileContent), &configuration)
			} else if filename == cfg.ProjectFile {
				return errors.New("Error: no such file or directory")
			}
			return errors.New("Wrong test file")
		},
		goToProjectPath: func(path string) error {
			if path == "/path2/" {
				return errors.New("Is ok")
			}
			return errors.New("Project project2 was not found ")
		},
		dirExists: func(path string) (bool, error) {
			return true, nil
		},
	}

	cfg.SetFileSystem(tfs)
	dialog := &dialog.Dialog{
		SelectProject: func(projects []string) (int, string, error) {
			return 0, "project2", nil
		},
		AddProjectName: func() (string, error) {
			return "", errors.New("Should not be called as the project name has been typed")
		},
		AddProjectPath: func(path string) (string, error) {
			return path, errors.New("Should not be called as the project has been selected")
		},
	}

	err := seekPath(cfg, dialog, tfs)
	assert.EqualError(t, err, "Is ok")
}

func TestSeekPath_runDialogCase4(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	tfs := &testFileSystem{
		getUserDirectory: func() (string, error) {
			return "/user/path/", nil
		},
		saveConfigFile: func(data interface{}, fileName string) error {
			return nil
		},
		getWd: func() (string, error) {
			return "", nil
		},
		readConfigFile: func(filename string, configuration interface{}) error {
			if filename == cfg.UserFile {
				return json.Unmarshal([]byte("{}"), &configuration)
			} else if filename == cfg.ProjectFile {
				return errors.New("Error: no such file or directory")
			}
			return errors.New("Wrong test file")
		},
		goToProjectPath: func(path string) error {
			if path == "/path2/" {
				return errors.New("Is ok")
			}
			return errors.New("Project project2 was not found ")
		},
		dirExists: func(path string) (bool, error) {
			return true, nil
		},
	}

	cfg.SetFileSystem(tfs)
	dialog := &dialog.Dialog{
		SelectProject: func(projects []string) (int, string, error) {
			return -1, "", nil
		},
		AddProjectName: func() (string, error) {
			return "", errors.New("Add project name dialog error")
		},
		AddProjectPath: func(path string) (string, error) {
			return path, errors.New("Should not be called as the project has been selected")
		},
	}

	err := seekPath(cfg, dialog, tfs)
	assert.EqualError(t, err, "Add project name dialog error")
}

func TestSeekPath_runDialogCase5(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	expectedName := "added project name"
	expectedPath := "/added/path/to/project"
	expectedContainer := ""

	var name, path, container string

	tfs := &testFileSystem{
		getUserDirectory: func() (string, error) {
			return "/user/path/", nil
		},
		saveConfigFile: func(data interface{}, fileName string) error {
			if fileName == "/added/path/to/project/"+cfg.ProjectFile {
				// expected data type conifg.ProjectConfig
				methodPath := reflect.ValueOf(data).MethodByName("GetPath")
				methodName := reflect.ValueOf(data).MethodByName("GetName")
				methodContainer := reflect.ValueOf(data).MethodByName("GetMainContainer")

				name = methodName.Call(nil)[0].Interface().(string)
				path = methodPath.Call(nil)[0].Interface().(string)
				container = methodContainer.Call(nil)[0].Interface().(string)

				return errors.New("Cannot save project config file")
			}

			return nil
		},
		getWd: func() (string, error) {
			return "", nil
		},
		readConfigFile: func(filename string, configuration interface{}) error {
			if filename == cfg.UserFile {
				return json.Unmarshal([]byte("{}"), &configuration)
			} else if filename == cfg.ProjectFile {
				return errors.New("Error: no such file or directory")
			}
			return errors.New("Wrong test file")
		},
		goToProjectPath: func(path string) error {
			if path == "/path2/" {
				return errors.New("Is ok")
			}
			return errors.New("Project project2 was not found ")
		},
		dirExists: func(path string) (bool, error) {
			return true, nil
		},
	}

	cfg.SetFileSystem(tfs)
	dialog := &dialog.Dialog{
		SelectProject: func(projects []string) (int, string, error) {
			return -1, "", nil
		},
		AddProjectName: func() (string, error) {
			return "added project name", nil
		},
		AddProjectPath: func(path string) (string, error) {
			return "/added/path/to/project", nil
		},
	}

	err := seekPath(cfg, dialog, tfs)

	assert.Equal(t, expectedName, name)
	assert.Equal(t, expectedPath, path)
	assert.Equal(t, expectedContainer, container)
	assert.EqualError(t, err, "Cannot save project config file")
}

// The same as TestSeekPath_runDialogCase5, but checks return nil for seek function
// see saveConfigFile return change
func TestSeekPath_runDialogCase6(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	expectedName := "added project name"
	expectedPath := "/added/path/to/project"
	expectedContainer := ""

	var name, path, container string

	tfs := &testFileSystem{
		getUserDirectory: func() (string, error) {
			return "/user/path/", nil
		},
		saveConfigFile: func(data interface{}, fileName string) error {
			if fileName == "/added/path/to/project/"+cfg.ProjectFile {
				// expected data type conifg.ProjectConfig
				methodPath := reflect.ValueOf(data).MethodByName("GetPath")
				methodName := reflect.ValueOf(data).MethodByName("GetName")
				methodContainer := reflect.ValueOf(data).MethodByName("GetMainContainer")

				name = methodName.Call(nil)[0].Interface().(string)
				path = methodPath.Call(nil)[0].Interface().(string)
				container = methodContainer.Call(nil)[0].Interface().(string)
			}

			return nil
		},
		getWd: func() (string, error) {
			return "", nil
		},
		readConfigFile: func(filename string, configuration interface{}) error {
			if filename == cfg.UserFile {
				return json.Unmarshal([]byte("{}"), &configuration)
			} else if filename == cfg.ProjectFile {
				return errors.New("Error: no such file or directory")
			}
			return errors.New("Wrong test file")
		},
		goToProjectPath: func(path string) error {
			if path == "/added/path/to/project" {
				return nil
			}
			return errors.New("Project project2 was not found ")
		},
		dirExists: func(path string) (bool, error) {
			return true, nil
		},
	}

	cfg.SetFileSystem(tfs)
	dialog := &dialog.Dialog{
		SelectProject: func(projects []string) (int, string, error) {
			return -1, "", nil
		},
		AddProjectName: func() (string, error) {
			return "added project name", nil
		},
		AddProjectPath: func(path string) (string, error) {
			return "/added/path/to/project", nil
		},
	}

	err := seekPath(cfg, dialog, tfs)

	assert.Equal(t, expectedName, name)
	assert.Equal(t, expectedPath, path)
	assert.Equal(t, expectedContainer, container)
	assert.Nil(t, err)
}
