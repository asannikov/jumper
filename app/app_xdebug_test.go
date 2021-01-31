package app

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/asannikov/jumper/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestXDebug(t *testing.T) {
	os.Args = []string{"jumper", "xdebug:fpm:enable"}

	c := &cli.App{}

	// @todo create mock object for dialog like testFileSystem
	DLG := testDialog{}
	DLG.setSelectProject = func(projects []string) (int, string, error) {
		return 0, "Project Name", nil
	}
	DLG.setAddProjectName = func() (string, error) {
		return "Project Name", nil
	}
	DLG.setAddProjectPath = func(path string) (string, error) {
		return "/current/path/", nil
	}
	DLG.setDockerService = func() (string, error) {
		return "start docker", nil
	}
	DLG.setMainContaner = func(cl []string) (int, string, error) {
		return 0, "container_name", nil
	}
	DLG.setDockerCliXdebugIniFilePath = func(p string) (string, error) {
		return "/path/to/cli/xdebug/ini", nil
	}
	DLG.setDockerFmpXdebugIniFilePath = func(p string) (string, error) {
		return "/path/to/Fmp/xdebug/ini", nil
	}
	DLG.setXdebugFileConfigLocation = func() (int, string, error) {
		return 0, "local", nil
	}

	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	jcfg := &jumperAppTest{
		dlg: &DLG,
		cfg: cfg,
		fs:  &testFileSystem{},
	}

	jcfg.fs.fileExists = func(file string) (bool, error) {
		return true, nil
	}

	jcfg.fs.getUserDirectory = func() (string, error) {
		return "/user/path/", nil
	}

	jcfg.fs.readConfigFile = func(filename string, configuration interface{}) error {
		return json.Unmarshal([]byte(testUserFileContent), &configuration)
	}

	jcfg.fs.getWd = func() (string, error) {
		return "/current/path/", nil
	}

	jcfg.fs.saveConfigFile = func(data interface{}, fileName string) error {
		return nil
	}

	jcfg.fs.goToProjectPath = func(path string) error {
		return nil
	}

	JumperAppTest(c, jcfg)

	err := c.Run(os.Args)

	assert.Nil(t, err)
}
