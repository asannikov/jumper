package command

import (
	"errors"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

type args struct {
	get   string
	slice []string
	tail  []string
}

func (a *args) Get(n int) string {
	return a.get
}

func (a *args) Slice() []string {
	return a.slice
}

func (a *args) First() string {
	return ""
}

func (a *args) Tail() []string {
	return a.tail
}

func (a *args) Len() int {
	return 0
}

func (a *args) Present() bool {
	return true
}

type testCli struct {
	args    map[string][]string
	command map[string]string
}

func (tc *testCli) GetCommand(cmd string, cfg commandHandleProjectConfig) string {
	return tc.command[cmd]
}
func (tc *testCli) GetArgs() map[string][]string {
	return tc.args
}

type testContainerlist struct {
	containerList []string
	err           error
}

func (tcl *testContainerlist) GetContainerList() ([]string, error) {
	return tcl.containerList, tcl.err
}

type testCliCommandHandleProjectConfig struct {
	mainContainer                    string
	saveContainerNameToProjectConfig error
	getShell                         string
}

func (tc *testCliCommandHandleProjectConfig) GetProjectMainContainer() string {
	return tc.mainContainer
}
func (tc *testCliCommandHandleProjectConfig) SaveContainerNameToProjectConfig(container string) error {
	return tc.saveContainerNameToProjectConfig
}
func (tc *testCliCommandHandleProjectConfig) GetShell() string {
	return tc.getShell
}

type testCallCliCommandDialog struct {
	setMainContaner func([]string) (int, string, error)
}

func (d *testCallCliCommandDialog) SetMainContaner(list []string) (int, string, error) {
	return d.setMainContaner(list)
}

type testCallCliCommandOptions struct {
	getExecCommand   func(ExecOptions, *cli.App) error
	getInitFunction  func(bool) string
	getContainerList func() ([]string, error)
}

func (x *testCallCliCommandOptions) GetExecCommand() func(ExecOptions, *cli.App) error {
	return x.getExecCommand
}
func (x *testCallCliCommandOptions) GetInitFunction() func(bool) string {
	return x.getInitFunction
}
func (x *testCallCliCommandOptions) GetContainerList() ([]string, error) {
	return x.getContainerList()
}

func TestCliHandleCase1(t *testing.T) {
	cfg := &testCliCommandHandleProjectConfig{
		mainContainer:                    "",
		saveContainerNameToProjectConfig: nil,
		getShell:                         "",
	}

	dlg := &testCallCliCommandDialog{
		setMainContaner: func(l []string) (int, string, error) {
			return 0, "", nil
		},
	}

	cli := &testCli{}

	a := &args{
		get:   "",
		slice: []string{},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{},
	}

	_, err := cliCommandHandle("cli", cfg, dlg, cli, cl, a)

	assert.EqualError(t, err, "Container name is empty. Set the container name")
}

func TestCliHandleCase2(t *testing.T) {
	cfg := &testCliCommandHandleProjectConfig{
		mainContainer:                    "containerName",
		saveContainerNameToProjectConfig: nil,
		getShell:                         "",
	}

	dlg := &testCallCliCommandDialog{
		setMainContaner: func(l []string) (int, string, error) {
			return 0, "", nil
		},
	}

	cli := &testCli{}

	a := &args{
		get:   "",
		slice: []string{},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{},
	}

	_, err := cliCommandHandle("cli", cfg, dlg, cli, cl, a)

	assert.EqualError(t, err, "Please specify a CLI command (ex. ls)")
}

func TestCliHandleCase3(t *testing.T) {
	cfg := &testCliCommandHandleProjectConfig{
		mainContainer:                    "containerName",
		saveContainerNameToProjectConfig: nil,
		getShell:                         "",
	}

	dlg := &testCallCliCommandDialog{
		setMainContaner: func(l []string) (int, string, error) {
			return 0, "", nil
		},
	}

	cli := &testCli{
		args: map[string][]string{
			"cli": []string{"-it"},
		},
		command: map[string]string{
			"cli": "sh",
		},
	}

	a := &args{
		get:   "",
		slice: []string{},
	}

	cl := &testContainerlist{
		err:           errors.New("GetContainerList error"),
		containerList: []string{},
	}

	args, err := cliCommandHandle("cli", cfg, dlg, cli, cl, a)

	assert.EqualError(t, err, "GetContainerList error")
	assert.Equal(t, []string{}, args)
}

func TestCliHandleCase4(t *testing.T) {
	cfg := &testCliCommandHandleProjectConfig{
		mainContainer:                    "containerName",
		saveContainerNameToProjectConfig: nil,
		getShell:                         "",
	}

	dlg := &testCallCliCommandDialog{
		setMainContaner: func(l []string) (int, string, error) {
			return 0, "", nil
		},
	}

	cli := &testCli{
		args: map[string][]string{
			"cli": []string{"-it"},
		},
		command: map[string]string{
			"cli": "sh",
		},
	}

	a := &args{
		get:   "",
		slice: []string{},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{},
	}

	args, err := cliCommandHandle("cli", cfg, dlg, cli, cl, a)

	assert.Nil(t, err)
	assert.Equal(t, []string{"exec", "-it", "containerName", "sh"}, args)
}

func TestCallCliCommandCase1(t *testing.T) {
	cfg := &testCliCommandHandleProjectConfig{
		mainContainer:                    "container_name",
		saveContainerNameToProjectConfig: nil,
		getShell:                         "",
	}

	dlg := &testCallCliCommandDialog{
		setMainContaner: func(l []string) (int, string, error) {
			return 0, "", nil
		},
	}

	opt := &testCallCliCommandOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return nil
		},
	}

	ctx := &cli.Context{
		App: &cli.App{},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{"ls"})

	ctx = cli.NewContext(&cli.App{}, set, ctx)

	app := CallCliCommand("cli", cfg, dlg, opt)
	assert.Nil(t, app.Action(ctx))
}

func TestCallCliCommandCase2(t *testing.T) {
	cfg := &testCliCommandHandleProjectConfig{
		mainContainer:                    "container_name",
		saveContainerNameToProjectConfig: nil,
		getShell:                         "",
	}

	dlg := &testCallCliCommandDialog{
		setMainContaner: func(l []string) (int, string, error) {
			return 0, "", nil
		},
	}

	opt := &testCallCliCommandOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return nil
		},
	}

	ctx := &cli.Context{
		App: &cli.App{},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx = cli.NewContext(&cli.App{}, set, ctx)

	app := CallCliCommand("cli", cfg, dlg, opt)
	assert.EqualError(t, app.Action(ctx), "Please specify a CLI command (ex. ls)")
}

func TestCallCliCommandCase3(t *testing.T) {
	cfg := &testCliCommandHandleProjectConfig{
		mainContainer:                    "",
		saveContainerNameToProjectConfig: nil,
		getShell:                         "",
	}

	dlg := &testCallCliCommandDialog{
		setMainContaner: func(l []string) (int, string, error) {
			return 0, "", nil
		},
	}

	opt := &testCallCliCommandOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return nil
		},
	}

	ctx := &cli.Context{
		App: &cli.App{},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx = cli.NewContext(&cli.App{}, set, ctx)

	app := CallCliCommand("cli", cfg, dlg, opt)
	assert.EqualError(t, app.Action(ctx), "Container name is empty. Set the container name")
}
