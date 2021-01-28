package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/asannikov/jumper/app/config"
	"github.com/asannikov/jumper/app/dialog"
	"github.com/asannikov/jumper/app/lib"
	"github.com/urfave/cli/v2"

	"github.com/stretchr/testify/assert"
)

const testProjectFileContent = `{
	"path": "/path/to/project/",
	"name": "project name",
	"main_container": ""
}`

const testUserFileContent = `{
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

type jumperAppTest struct {
	dlg *dialog.Dialog
	cfg *config.Config
	fs  *FileSystem
}

func JumperAppTest(cli *cli.App, jat *jumperAppTest) {
	jat.cfg.Init()
	jat.cfg.SetFileSystem(jat.fs)

	// Loading only global config
	loadGlobalConfig(jat.cfg, jat.fs)

	// Loading project config if exists
	loadProjectConfig(jat.cfg, jat.fs)

	// Define docker command
	defineDockerCommand(jat.cfg, jat.dlg)

	initf := func(seekProject bool) string {
		if err := seekPath(jat.cfg, jat.dlg, jat.fs, seekProject); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if seekProject == true {
			currentDir, _ := jat.fs.GetWd()
			fmt.Printf("\nchanged user location to directory: %s\n\n", currentDir)
			return currentDir
		}

		return ""
	}

	cli.Copyright = lib.GetCopyrightText(jat.cfg)
	cli.Commands = commandList(jat.cfg, jat.dlg, initf)
}

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

// TestDefinePaths tests definePaths
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
	assert.Equal(t, "/user/path/"+string(os.PathSeparator)+".jumper.json", cfg.GetUserFile())

	tfs = &testFileSystem{
		getUserDirectory: func() (string, error) {
			return "", errors.New("Cannot read directory")
		},
	}

	err = definePaths(cfg, tfs)
	assert.EqualError(t, err, "Cannot read directory")
}

// TestSeekPathDefinePaths tests seekPath
func TestSeekPathDefinePaths(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	tfs := &testFileSystem{
		getUserDirectory: func() (string, error) {
			return "", errors.New("Cannot read directory")
		},
	}

	dialog := &dialog.Dialog{}

	err := seekPath(cfg, dialog, tfs, true)
	assert.EqualError(t, err, "Cannot read directory")
}

// TestSeekPathLoadConfig tests seekPath
func TestSeekPathLoadConfig(t *testing.T) {
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

	err := seekPath(cfg, dialog, tfs, true)
	assert.EqualError(t, err, "Cannot read config file")
}

// TestSeekPathGetWd tests seekPath
func TestSeekPathGetWd(t *testing.T) {
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

	err := seekPath(cfg, dialog, tfs, true)
	assert.EqualError(t, err, "Cannot get current directory")
}

// TestSeekPathGetProjectNameList tests seekPath
func TestSeekPathGetProjectNameList(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	cfg.Init()

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
			if filename == cfg.GetUserFile() {
				return json.Unmarshal([]byte(testUserFileContent), &configuration)
			} else if filename == cfg.ProjectFile {
				return json.Unmarshal([]byte(testProjectFileContent), &configuration)
			}
			return errors.New("Wrong test file")
		},
	}

	cfg.SetFileSystem(tfs)
	dialog := &dialog.Dialog{}

	err := seekPath(cfg, dialog, tfs, true)
	assert.EqualError(t, err, "stop execution")

	pl := cfg.GetProjectNameList()
	assert.Equal(t, 2, len(pl))
}

// TestSeekPathrunDialogCase1 tests seekPath
func TestSeekPathrunDialogCase1(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	cfg.Init()

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
			if filename == cfg.GetUserFile() {
				return json.Unmarshal([]byte(testUserFileContent), &configuration)
			} else if filename == cfg.ProjectFile {
				return errors.New("Error: no such file or directory")
			}
			return errors.New("Wrong test file")
		},
	}

	cfg.SetFileSystem(tfs)
	dialog := &dialog.Dialog{}
	dialog.SetSelectProjectTest(func(projects []string) (int, string, error) {
		return 0, "", errors.New("Error SelectProject dialog")
	})

	err := seekPath(cfg, dialog, tfs, true)
	assert.EqualError(t, err, "Error SelectProject dialog")
}

// TestSeekPathrunDialogCase2 tests seekPath
func TestSeekPathrunDialogCase2(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	cfg.Init()

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
			if filename == cfg.GetUserFile() {
				return json.Unmarshal([]byte(testUserFileContent), &configuration)
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
	dialog := &dialog.Dialog{}
	dialog.SetSelectProjectTest(func(projects []string) (int, string, error) {
		return -1, "New Project Name", nil
	})
	dialog.SetAddProjectNameTest(func() (string, error) {
		return "", errors.New("Should not be called as the project name has been typed")
	})
	dialog.SetAddProjectPathTest(func(path string) (string, error) {
		if path == "/current/path/" {
			return path, errors.New("Is ok")
		}
		return "", errors.New("Current path was not set")
	})

	err := seekPath(cfg, dialog, tfs, true)
	assert.EqualError(t, err, "Is ok")
}

// TestSeekPathrunDialogCase3 tests seekPath
func TestSeekPathrunDialogCase3(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	cfg.Init()

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
			if filename == cfg.GetUserFile() {
				return json.Unmarshal([]byte(testUserFileContent), &configuration)
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
	dialog := &dialog.Dialog{}
	dialog.SetSelectProjectTest(func(projects []string) (int, string, error) {
		return 0, "project2", nil
	})
	dialog.SetAddProjectNameTest(func() (string, error) {
		return "", errors.New("Should not be called as the project name has been typed")
	})
	dialog.SetAddProjectPathTest(func(path string) (string, error) {
		return path, errors.New("Should not be called as the project has been selected")
	})

	err := seekPath(cfg, dialog, tfs, true)
	assert.EqualError(t, err, "Is ok")
}

// TestSeekPathrunDialogCase4 tests seekPath
func TestSeekPathrunDialogCase4(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	cfg.Init()

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
			if filename == cfg.GetUserFile() {
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
	dialog := &dialog.Dialog{}
	dialog.SetSelectProjectTest(func(projects []string) (int, string, error) {
		return -1, "", nil
	})
	dialog.SetAddProjectNameTest(func() (string, error) {
		return "", errors.New("Add project name dialog error")
	})
	dialog.SetAddProjectPathTest(func(path string) (string, error) {
		return path, errors.New("Should not be called as the project has been selected")
	})

	err := seekPath(cfg, dialog, tfs, true)
	assert.EqualError(t, err, "Add project name dialog error")
}

// TestSeekPathrunDialogCase5 tests seekPath
func TestSeekPathrunDialogCase5(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	cfg.Init()

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
			if filename == cfg.GetUserFile() {
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
	dialog := &dialog.Dialog{}
	dialog.SetSelectProjectTest(func(projects []string) (int, string, error) {
		return -1, "", nil
	})
	dialog.SetAddProjectNameTest(func() (string, error) {
		return "added project name", nil
	})
	dialog.SetAddProjectPathTest(func(path string) (string, error) {
		return "/added/path/to/project", nil
	})

	err := seekPath(cfg, dialog, tfs, true)

	assert.Equal(t, expectedName, name)
	assert.Equal(t, expectedPath, path)
	assert.Equal(t, expectedContainer, container)
	assert.EqualError(t, err, "Cannot save project config file")
}

// TestSeekPathrunDialogCase6 tests seekPath. The same as TestSeekPathrunDialogCase5, but checks return nil for seek function
// see saveConfigFile return change
func TestSeekPathrunDialogCase6(t *testing.T) {
	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	cfg.Init()

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
			if filename == cfg.GetUserFile() {
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

	dialog := &dialog.Dialog{}
	dialog.SetSelectProjectTest(func(projects []string) (int, string, error) {
		return -1, "", nil
	})
	dialog.SetAddProjectNameTest(func() (string, error) {
		return "added project name", nil
	})
	dialog.SetAddProjectPathTest(func(path string) (string, error) {
		return "/added/path/to/project", nil
	})

	err := seekPath(cfg, dialog, tfs, true)

	assert.Equal(t, expectedName, name)
	assert.Equal(t, expectedPath, path)
	assert.Equal(t, expectedContainer, container)
	assert.Nil(t, err)
}
