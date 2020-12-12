package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// FileSystem file management type
type FileSystem struct{}

// FileExists checks if file exists
func (fs *FileSystem) FileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, err
	}
	return !info.IsDir(), nil
}

// DirExists checks if directory exists
func (fs *FileSystem) DirExists(path string) (bool, error) {
	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			return true, nil
		}
		return false, fmt.Errorf("Path %s is a file ", path)
	} else if os.IsNotExist(err) {
		// path does *not* exist
		return false, err
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false, err
	}
}

// ReadConfigFile reads json file
func (fs *FileSystem) ReadConfigFile(filename string, configuration interface{}) (err error) {
	if len(filename) == 0 {
		return fmt.Errorf("Reading config error, file name is empty")
	}

	if t, err := fs.FileExists(filename); err != nil && t == false {
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), &configuration)
}

// SaveConfigFile save data to json file
func (fs *FileSystem) SaveConfigFile(data interface{}, fileName string) error {
	file, _ := json.MarshalIndent(data, "", " ")
	err := ioutil.WriteFile(fileName, file, 0644)
	return err
}

// GoToProjectPath does CD to the path
func (fs *FileSystem) GoToProjectPath(projectPath string) error {
	projectPath = strings.TrimRight(projectPath, string(os.PathSeparator))
	projectPath, err := filepath.EvalSymlinks(projectPath)

	if err != nil {
		return err
	}

	if err := os.Chdir(projectPath); err != nil {
		return err
	}

	currentDir, err := fs.GetWd()

	if err != nil {
		return err
	}

	if currentDir != projectPath {
		return fmt.Errorf("Expected path %s, the current one %s", projectPath, currentDir)
	}

	return nil
}

// GetWd works as `pwd`
func (fs *FileSystem) GetWd() (string, error) {
	return os.Getwd()
}

// GetUserDirectory gets home user's directory
func (fs *FileSystem) GetUserDirectory() (string, error) {
	usr, err := user.Current()

	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
