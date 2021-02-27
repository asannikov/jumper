package command

import (
	"errors"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

type testMagentoGlobalConfig struct {
	saveContainerNameToProjectConfig error
	getProjectMainContainer          string
	saveDockerProjectPath            error
	getProjectDockerPath             string
}

func (c *testMagentoGlobalConfig) SaveContainerNameToProjectConfig(string) error {
	return c.saveContainerNameToProjectConfig
}
func (c *testMagentoGlobalConfig) GetProjectMainContainer() string {
	return c.getProjectMainContainer
}
func (c *testMagentoGlobalConfig) SaveDockerProjectPath(string) error {
	return c.saveDockerProjectPath
}
func (c *testMagentoGlobalConfig) GetProjectDockerPath() string {
	return c.getProjectDockerPath
}

type testMagentoDialog struct {
	setMainContaner   func([]string) (int, string, error)
	dockerProjectPath func(string) (string, error)
}

func (d *testMagentoDialog) SetMainContaner(l []string) (int, string, error) {
	return d.setMainContaner(l)
}
func (d *testMagentoDialog) DockerProjectPath(p string) (string, error) {
	return d.dockerProjectPath(p)
}

type testMagentoOptions struct {
	getExecCommand     func(ExecOptions, *cli.App) error
	getInitFunction    func(bool) string
	getCommandLocation func(string, string) (string, error)
	getContainerList   func() ([]string, error)
	checkMagentoBin    func(string, string) (bool, error)
}

func (o *testMagentoOptions) GetExecCommand() func(eo ExecOptions, a *cli.App) error {
	return o.getExecCommand
}
func (o *testMagentoOptions) GetCommandLocation() func(string, string) (string, error) {
	return o.getCommandLocation
}
func (o *testMagentoOptions) GetInitFunction() func(bool) string {
	return o.getInitFunction
}
func (o *testMagentoOptions) GetContainerList() ([]string, error) {
	return o.getContainerList()
}

func (o *testMagentoOptions) CheckMagentoBin(containerName string, magentoBin string) (bool, error) {
	return o.checkMagentoBin(containerName, magentoBin)
}

func TestCallMagentoCommandCase1(t *testing.T) {
	cfg := &testMagentoGlobalConfig{}
	dlg := &testMagentoDialog{}
	opt := &testMagentoOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, errors.New("getContainerList list error")
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := callMagentoCommanBin(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "getContainerList list error")
}
