package app

import (
	"os"
	"strings"
	"testing"

	"jumper/app/config" // github.com/asannikov/

	"github.com/stretchr/testify/assert"
)

func testFileExistsCase1(t *testing.T) {
	fs := FileSystem{}
	found, err := fs.DirExists("/test/folder")
	assert.Equal(t, false, found)
	assert.EqualError(t, err, "stat /test/folder: no such file or directory")
}

func testFileExistsCase2(t *testing.T) {
	fs := FileSystem{}
	currentDir, _ := fs.GetWd()
	found, err := fs.DirExists(currentDir)
	assert.True(t, found)
	assert.Equal(t, nil, err)
}

func testFileExistsCase3(t *testing.T) {
	fs := FileSystem{}
	currentDir, _ := fs.GetWd()
	found, err := fs.DirExists(currentDir + string(os.PathSeparator) + "main_file_test.go")
	assert.False(t, found)
	assert.True(t, strings.Contains(err.Error(), "main_file_test.go is a file"))
}

func testReadConfigFileCase1(t *testing.T) {
	fs := FileSystem{}
	configuration := config.ProjectConfig{}
	err := fs.ReadConfigFile("", &configuration)
	assert.EqualError(t, err, "Reading config error, file name is empty")
}

func testReadConfigFile2(t *testing.T) {
	fs := FileSystem{}
	currentDir, _ := fs.GetWd()
	configuration := config.ProjectConfig{}
	err := fs.ReadConfigFile(currentDir+string(os.PathSeparator)+"jumper_test.json", &configuration)
	assert.Equal(t, nil, err)
	assert.Equal(t, "main_container_name", configuration.GetMainContainer())
}

func testGoToProjectPath1(t *testing.T) {
	fs := FileSystem{}
	currentDir, _ := fs.GetWd()
	// currentDir := currentDir + string(os.PathSeparator) + "config" // for local testing
	err := fs.GoToProjectPath(currentDir)
	assert.Equal(t, nil, err)
}

func testGoToProjectPath2(t *testing.T) {
	fs := FileSystem{}
	currentDir, _ := fs.GetWd()
	path := currentDir + string(os.PathSeparator) + "nofolder"
	err := fs.GoToProjectPath(path)
	assert.True(t, strings.Contains(err.Error(), "no such file or directory"))
}
