package app

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/asannikov/jumper/app/command"
	"github.com/asannikov/jumper/app/config"
	"github.com/asannikov/jumper/app/lib"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

type jumperAppTest struct {
	dlg *testDialog
	cfg *config.Config
	fs  *testFileSystem
}

func getTestDialog() *testDialog {
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

	return &DLG
}

func prepareTestFs(jcfg *jumperAppTest) {
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
}

func getTestDockerInstance() *testDockerInstance {
	dck := &testDockerInstance{}
	dck.stat = func() (string, error) {
		return "", nil
	}
	dck.initClient = func() error {
		return nil
	}
	dck.getContainerList = func() ([]string, error) {
		return []string{"container"}, nil
	}
	dck.getClient = func() *client.Client {
		return &client.Client{}
	}

	return dck
}

func jumperMainAppTest() (*cli.App, *jumperAppTest, *commandOptions) {
	cliApp := &cli.App{}

	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	jcfg := &jumperAppTest{
		dlg: getTestDialog(),
		cfg: cfg,
		fs:  &testFileSystem{},
	}

	prepareTestFs(jcfg)

	jcfg.cfg.Init()
	jcfg.cfg.SetFileSystem(jcfg.fs)

	// Loading only global config
	loadGlobalConfig(jcfg.cfg, jcfg.fs)

	// Loading project config if exists
	loadProjectConfig(jcfg.cfg, jcfg.fs)

	// Define docker command
	defineDockerCommand(jcfg.cfg, jcfg.dlg)

	opt := getOptions(jcfg.cfg, jcfg.dlg)
	opt.setInitFuntion(func(seekProject bool) string {
		if err := seekPath(jcfg.cfg, jcfg.dlg, jcfg.fs, seekProject); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if seekProject == true {
			currentDir, _ := jcfg.fs.GetWd()
			fmt.Printf("\nchanged user location to directory: %s\n\n", currentDir)
			return currentDir
		}

		return ""
	})

	dockerDialog := getDockerStartDialog()
	dockerDialog.setDialog(jcfg.dlg)
	dockerDialog.setDocker(getTestDockerInstance())
	dockerDialog.setDockerService(jcfg.cfg.GetDockerCommand())
	dockerDialog.startDockerDialog = func(cl *dockerStartDialog) (string, error) {
		return "", nil
	}
	dockerDialog.containerList = func(cl *dockerStartDialog) ([]string, error) {
		return []string{"container"}, nil
	}
	opt.setDockerDialog(dockerDialog)

	cliApp.Copyright = lib.GetCopyrightText(jcfg.cfg)

	return cliApp, jcfg, opt
}

func TestAppCall(t *testing.T) {
	cliApp, jcfg, opt := jumperMainAppTest()
	opt.execCommand = func(eo command.ExecOptions, c *cli.App) error {
		fmt.Printf("\ncommand: %s\n\n", eo.GetCommand()+" "+strings.Join(eo.GetArgs(), " "))
		return nil
	}

	cliApp.Commands = commandList(jcfg.cfg, jcfg.dlg, opt)
	assert.Nil(t, cliApp.Run([]string{"jumper", "xdebug:fpm:enable"}))
}
