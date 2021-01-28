package app

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/asannikov/jumper/app/config"
	"github.com/asannikov/jumper/app/dialog"
	"github.com/urfave/cli/v2"
)

func TestXDebug(t *testing.T) {
	os.Args = []string{"jumper", "xdebug:fpm:enable"}

	c := &cli.App{}

	DLG := dialog.InitDialogFunctions()
	DLG.SetSelectProjectTest(func(projects []string) (int, string, error) {
		return 0, "Project Name", nil
	})
	DLG.SetAddProjectNameTest(func() (string, error) {
		return "Project Name", nil
	})
	DLG.SetAddProjectPathTest(func(path string) (string, error) {
		return "/current/path/", nil
	})

	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	jcfg := &jumperAppTest{
		dlg: &DLG,
		cfg: cfg,
		fs:  &testFileSystem{},
	}

	jcfg.fs.fileExists = func(file string) (bool, error) {
		log.Println(file)
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

	log.Println(err)
	//assert.Nil(t, err)
}
