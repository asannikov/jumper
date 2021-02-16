package command

import (
	"errors"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

type testStopOptions struct {
	getExecCommand    func(ExecOptions, *cli.App) error
	getInitFunction   func(bool) string
	getContainerList  func() ([]string, error)
	getDockerStatus   bool
	getStopContainers func([]string) error
}

func (x *testStopOptions) GetInitFunction() func(bool) string {
	return x.getInitFunction
}
func (x *testStopOptions) GetContainerList() ([]string, error) {
	return x.getContainerList()
}
func (x *testStopOptions) GetExecCommand() func(ExecOptions, *cli.App) error {
	return x.getExecCommand
}
func (x *testStopOptions) GetDockerStatus() bool {
	return x.getDockerStatus
}
func (x *testStopOptions) GetStopContainers() func([]string) error {
	return x.getStopContainers
}

type testStopConfig struct {
	getProjectMainContainer          string
	saveContainerNameToProjectConfig error
}

func (tsc *testStopConfig) GetProjectMainContainer() string {
	return tsc.getProjectMainContainer
}
func (tsc *testStopConfig) SaveContainerNameToProjectConfig(name string) error {
	return tsc.saveContainerNameToProjectConfig
}

type testStopDialog struct {
	setMainContaner func([]string) (int, string, error)
}

func (tsd *testStopDialog) SetMainContaner(list []string) (int, string, error) {
	return tsd.setMainContaner(list)
}

func TestCallStopAllContainersCommandCase1(t *testing.T) {
	opt := &testStopOptions{
		getDockerStatus: false,
		getInitFunction: func(s bool) string {
			return ""
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStopAllContainersCommand(opt)

	assert.EqualError(t, app.Action(ctx), "Docker is not running")
}

func TestCallStopAllContainersCommandCase2(t *testing.T) {
	opt := &testStopOptions{
		getDockerStatus: true,
		getInitFunction: func(s bool) string {
			return ""
		},
		getStopContainers: func(list []string) error {
			return nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStopAllContainersCommand(opt)

	assert.Nil(t, app.Action(ctx))
}

func TestCallStopMainContainerCommandCase1(t *testing.T) {
	opt := &testStopOptions{
		getDockerStatus: false,
	}

	cfg := &testStopConfig{}
	dlg := &testStopDialog{}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStopMainContainerCommand(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "Docker is not running")
}

func TestCallStopMainContainerCommandCase2(t *testing.T) {
	opt := &testStopOptions{
		getDockerStatus: true,
		getContainerList: func() ([]string, error) {
			return []string{}, errors.New("get container error")
		},
		getInitFunction: func(s bool) string {
			return ""
		},
	}

	cfg := &testStopConfig{}
	dlg := &testStopDialog{}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStopMainContainerCommand(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "get container error")
}

func TestCallStopMainContainerCommandCase3(t *testing.T) {
	opt := &testStopOptions{
		getDockerStatus: true,
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
		getInitFunction: func(s bool) string {
			return ""
		},
	}

	cfg := &testStopConfig{
		getProjectMainContainer: "",
	}

	dlg := &testStopDialog{
		setMainContaner: func(list []string) (int, string, error) {
			return 0, "", errors.New("defineProjectMainContainer error")
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStopMainContainerCommand(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "defineProjectMainContainer error")
}

func TestCallStopMainContainerCommandCase4(t *testing.T) {
	opt := &testStopOptions{
		getDockerStatus: true,
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
		getInitFunction: func(s bool) string {
			return ""
		},
		getStopContainers: func(list []string) error {
			return errors.New("stopContainers error")
		},
	}

	cfg := &testStopConfig{
		getProjectMainContainer: "container_name",
	}

	dlg := &testStopDialog{}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStopMainContainerCommand(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "stopContainers error")
}
