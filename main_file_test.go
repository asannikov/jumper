package main

import (
	"mgt/config"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExistsCase1(t *testing.T) {
	fs := FileSystem{}
	found, err := fs.DirExists("/test/folder")
	assert.Equal(t, false, found)
	assert.EqualError(t, err, "stat /test/folder: no such file or directory")
}

func TestFileExistsCase2(t *testing.T) {
	fs := FileSystem{}
	currentDir, _ := fs.GetWd()
	found, err := fs.DirExists(currentDir)
	assert.True(t, found)
	assert.Equal(t, nil, err)
}

func TestFileExistsCase3(t *testing.T) {
	fs := FileSystem{}
	currentDir, _ := fs.GetWd()
	found, err := fs.DirExists(currentDir + string(os.PathSeparator) + "main_file_test.go")
	assert.False(t, found)
	assert.True(t, strings.Contains(err.Error(), "main_file_test.go is a file"))
}

func TestReadConfigFileCase1(t *testing.T) {
	fs := FileSystem{}
	configuration := config.ProjectConfig{}
	err := fs.ReadConfigFile("", &configuration)
	assert.EqualError(t, err, "Reading config error, file name is empty")
}

func TestReadConfigFile2(t *testing.T) {
	fs := FileSystem{}
	currentDir, _ := fs.GetWd()
	configuration := config.ProjectConfig{}
	err := fs.ReadConfigFile(currentDir+string(os.PathSeparator)+"mgt_test.json", &configuration)
	assert.Equal(t, nil, err)
	assert.Equal(t, "main_container_name", configuration.GetMainContainer())
}

func TestGoToProjectPath1(t *testing.T) {
	fs := FileSystem{}
	currentDir, _ := fs.GetWd()
	path := currentDir + string(os.PathSeparator) + "config"
	err := fs.GoToProjectPath(path)
	assert.Equal(t, nil, err)
}

func TestGoToProjectPath2(t *testing.T) {
	fs := FileSystem{}
	currentDir, _ := fs.GetWd()
	path := currentDir + string(os.PathSeparator) + "nofolder"
	err := fs.GoToProjectPath(path)
	assert.True(t, strings.Contains(err.Error(), "no such file or directory"))
}
