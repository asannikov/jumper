package command

import (
	"errors"
	"flag"
	"fmt"
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
func (c *testMagentoGlobalConfig) GetCommandInactveStatus(string) bool {
	return true
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
			return []string{"container"}, errors.New("getContainerList list error")
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

func TestCallMagentoCommandCase2(t *testing.T) {
	cfg := &testMagentoGlobalConfig{}
	dlg := &testMagentoDialog{
		setMainContaner: func([]string) (int, string, error) {
			return 0, "", errors.New("setMainContaner error")
		},
	}
	opt := &testMagentoOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := callMagentoCommanBin(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "setMainContaner error")
}

func TestCallMagentoCommandCase3(t *testing.T) {
	cfg := &testMagentoGlobalConfig{}
	dlg := &testMagentoDialog{
		setMainContaner: func([]string) (int, string, error) {
			return 0, "container_name", nil
		},
		dockerProjectPath: func(string) (string, error) {
			return "", errors.New("dockerProjectPath error")
		},
	}
	opt := &testMagentoOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := callMagentoCommanBin(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "dockerProjectPath error")
}

func TestCallMagentoCommandCase4(t *testing.T) {
	cfg := &testMagentoGlobalConfig{
		getProjectDockerPath:    "/var/www",
		getProjectMainContainer: "container_name",
	}
	dlg := &testMagentoDialog{
		setMainContaner: func([]string) (int, string, error) {
			return 0, "container_name", nil
		},
		dockerProjectPath: func(string) (string, error) {
			return "/var/www", nil
		},
	}
	opt := &testMagentoOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		checkMagentoBin: func(mainContiner string, path string) (bool, error) {
			assert.Equal(t, mainContiner, "container_name")
			assert.Equal(t, path, "/var/www/bin/magento")
			return false, errors.New("path check error")
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := callMagentoCommanBin(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "path check error")
}

func TestCallMagentoCommandCase5(t *testing.T) {
	cfg := &testMagentoGlobalConfig{
		getProjectDockerPath:    "/var/www/",
		getProjectMainContainer: "container_name",
	}
	dlg := &testMagentoDialog{
		setMainContaner: func([]string) (int, string, error) {
			return 0, "container_name", nil
		},
		dockerProjectPath: func(string) (string, error) {
			return "/var/www/", nil
		},
	}
	opt := &testMagentoOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		checkMagentoBin: func(mainContiner string, path string) (bool, error) {
			assert.Equal(t, mainContiner, "container_name")
			assert.Equal(t, path, "/var/www/bin/magento")
			return false, errors.New("path check error")
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := callMagentoCommanBin(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "path check error")
}

func TestCallMagentoCommandCase6(t *testing.T) {
	cfg := &testMagentoGlobalConfig{
		getProjectDockerPath:    "/var/www/",
		getProjectMainContainer: "container_name",
	}
	dlg := &testMagentoDialog{
		setMainContaner: func([]string) (int, string, error) {
			return 0, "container_name", nil
		},
		dockerProjectPath: func(string) (string, error) {
			return "/var/www/", nil
		},
	}
	opt := &testMagentoOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		checkMagentoBin: func(mainContiner string, path string) (bool, error) {
			return false, nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := callMagentoCommanBin(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), fmt.Sprintf("Cannot find magento root folder. Searched for: %s", []string{
		"bin/magento",
		"html/bin/magento",
		"source/bin/magento",
		"src/bin/magento",
	}))
}

func TestCallMagentoCommandCase7(t *testing.T) {
	cfg := &testMagentoGlobalConfig{
		getProjectDockerPath:    "/var/www/",
		getProjectMainContainer: "container_name",
	}
	dlg := &testMagentoDialog{
		setMainContaner: func([]string) (int, string, error) {
			return 0, "container_name", nil
		},
		dockerProjectPath: func(string) (string, error) {
			return "/var/www/", nil
		},
	}
	opt := &testMagentoOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		checkMagentoBin: func(mainContiner string, path string) (bool, error) {
			if path == "/var/www/source/bin/magento" {
				return true, nil
			}
			return false, nil
		},
		getExecCommand: func(ex ExecOptions, a *cli.App) error {
			assert.Equal(t, ex.GetArgs(), []string{
				"exec",
				"-it",
				"container_name",
				"/var/www/source/bin/magento",
			})
			return errors.New("Exec command error")
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := callMagentoCommanBin(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "Exec command error")
}

func TestCallMagentoCommandMageRunCase1(t *testing.T) {
	cfg := &testMagentoGlobalConfig{}
	dlg := &testMagentoDialog{}
	opt := &testMagentoOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, errors.New("getContainerList list error")
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := callMagentoCommandMageRun(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "getContainerList list error")
}

func TestCallMagentoCommandMageRunCase2(t *testing.T) {
	cfg := &testMagentoGlobalConfig{}
	dlg := &testMagentoDialog{
		setMainContaner: func([]string) (int, string, error) {
			return 0, "", errors.New("setMainContaner error")
		},
	}
	opt := &testMagentoOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := callMagentoCommandMageRun(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "setMainContaner error")
}

func TestCallMagentoCommandMageRunCase3(t *testing.T) {
	cfg := &testMagentoGlobalConfig{}
	dlg := &testMagentoDialog{
		setMainContaner: func([]string) (int, string, error) {
			return 0, "container_name", nil
		},
		dockerProjectPath: func(string) (string, error) {
			return "", errors.New("dockerProjectPath error")
		},
	}
	opt := &testMagentoOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		getCommandLocation: func(string, string) (string, error) {
			return "", errors.New("commandLocation error")
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := callMagentoCommandMageRun(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "commandLocation error")
}

func TestCallMagentoCommandMageRunCase4(t *testing.T) {
	cfg := &testMagentoGlobalConfig{
		getProjectDockerPath:    "/var/www",
		getProjectMainContainer: "container_name",
	}
	dlg := &testMagentoDialog{
		setMainContaner: func([]string) (int, string, error) {
			return 0, "container_name", nil
		},
		dockerProjectPath: func(string) (string, error) {
			return "/var/www", nil
		},
	}
	opt := &testMagentoOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		getCommandLocation: func(containerName string, mr string) (string, error) {
			return "/var/usr/bin/" + mr, nil
		},
		getExecCommand: func(ex ExecOptions, a *cli.App) error {
			assert.Equal(t, ex.GetArgs(), []string{
				"exec",
				"-it",
				"container_name",
				"/var/usr/bin/n98-magerun2.phar",
			})
			return errors.New("Exec command error")
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := callMagentoCommandMageRun(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "Exec command error")
}
