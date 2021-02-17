package command

import (
	"errors"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

type testStartOptions struct {
	getExecCommand    func(ExecOptions, *cli.App) error
	getInitFunction   func(bool) string
	getContainerList  func() ([]string, error)
	getDockerStatus   bool
	getStopContainers func([]string) error
}

func (x *testStartOptions) GetInitFunction() func(bool) string {
	return x.getInitFunction
}
func (x *testStartOptions) GetContainerList() ([]string, error) {
	return x.getContainerList()
}
func (x *testStartOptions) GetExecCommand() func(ExecOptions, *cli.App) error {
	return x.getExecCommand
}
func (x *testStartOptions) GetDockerStatus() bool {
	return x.getDockerStatus
}
func (x *testStartOptions) GetStopContainers() func([]string) error {
	return x.getStopContainers
}

type testStartConfig struct {
	getStartCommand                 string
	saveStartCommandToProjectConfig error
	projectMainContainer            string
}

func (c *testStartConfig) GetProjectMainContainer() string {
	return c.projectMainContainer
}

func (c *testStartConfig) GetStartCommand() string {
	return c.getStartCommand
}

func (c *testStartConfig) SaveStartCommandToProjectConfig(s string) error {
	return c.saveStartCommandToProjectConfig
}

func (c *testStartConfig) SaveContainerNameToProjectConfig(s string) error {
	return nil
}

type testStartDialog struct {
	startCommand    func() (string, error)
	setMainContaner func() (int, string, error)
}

func (d *testStartDialog) StartCommand() (string, error) {
	return d.startCommand()
}
func (d *testStartDialog) SetMainContaner([]string) (int, string, error) {
	return d.setMainContaner()
}

func TestDefineStartCommandCase1(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "start_comand",
	}

	dlg := &testStartDialog{}

	assert.Nil(t, defineStartCommand(cfg, dlg, []string{}))
}

func TestDefineStartCommandCase2(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "",
	}

	dlg := &testStartDialog{
		startCommand: func() (string, error) {
			return "start_command", errors.New("Start command error")
		},
	}

	assert.EqualError(t, defineStartCommand(cfg, dlg, []string{}), "Start command error")
}

func TestDefineStartCommandCase3(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "",
	}

	dlg := &testStartDialog{
		startCommand: func() (string, error) {
			return "", nil
		},
	}

	assert.EqualError(t, defineStartCommand(cfg, dlg, []string{}), "Start command cannot be empty")
}

func TestDefineStartCommandCase4(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand:                 "",
		saveStartCommandToProjectConfig: errors.New("saveStartCommandToProjectConfig error"),
	}

	dlg := &testStartDialog{
		startCommand: func() (string, error) {
			return "start_command", nil
		},
	}

	assert.EqualError(t, defineStartCommand(cfg, dlg, []string{}), "saveStartCommandToProjectConfig error")
}

func TestRunStartProjectCase1(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "",
	}

	opt := &testStartOptions{
		getExecCommand: func(ExecOptions, *cli.App) error {
			return nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{
		"container_name",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)

	assert.Nil(t, runStartProject(ctx, cfg, []string{}, opt))
}

func TestRunStartProjectCase2(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "start_command up",
	}

	opt := &testStartOptions{
		getExecCommand: func(e ExecOptions, c *cli.App) error {
			assert.Equal(t, e.GetCommand(), "start_command")
			assert.Equal(t, e.GetArgs(), []string{"up", "--force-recreate", "container_name"})
			return nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{
		"container_name",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)

	assert.Nil(t, runStartProject(ctx, cfg, []string{"--force-recreate"}, opt))
}

func TestCallStartProjectBasicCase1(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "start_command up",
	}

	dlg := &testStartDialog{}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, errors.New("GetContainerList error")
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectBasic(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "GetContainerList error")
}

func TestCallStartProjectBasicCase2(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "start_command up",
	}

	dlg := &testStartDialog{
		setMainContaner: func() (int, string, error) {
			return 0, "", nil
		},
	}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectBasic(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "Container name is empty. Set the container name")
}

func TestCallStartProjectBasicCase3(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "",
	}

	dlg := &testStartDialog{
		setMainContaner: func() (int, string, error) {
			return 0, "container_name", nil
		},
		startCommand: func() (string, error) {
			return "start_command", errors.New("Start command error")
		},
	}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
		getExecCommand: func(ExecOptions, *cli.App) error {
			return nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectBasic(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "Start command error")
}

func TestCallStartProjectBasicCase4(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "commandName",
	}

	dlg := &testStartDialog{
		setMainContaner: func() (int, string, error) {
			return 0, "container_name", nil
		},
		startCommand: func() (string, error) {
			return "start_command", errors.New("Start command error")
		},
	}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
		getExecCommand: func(ExecOptions, *cli.App) error {
			return nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectBasic(cfg, dlg, opt)

	assert.Nil(t, app.Action(ctx))
}

func TestCallStartProjectForceRecreateCase1(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "start_command up",
	}

	dlg := &testStartDialog{}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, errors.New("GetContainerList error")
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectForceRecreate(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "GetContainerList error")
}

func TestCallStartProjectForceRecreateCase2(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "start_command up",
	}

	dlg := &testStartDialog{
		setMainContaner: func() (int, string, error) {
			return 0, "", nil
		},
	}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectForceRecreate(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "Container name is empty. Set the container name")
}

func TestCallStartProjectForceRecreateCase3(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "",
	}

	dlg := &testStartDialog{
		setMainContaner: func() (int, string, error) {
			return 0, "container_name", nil
		},
		startCommand: func() (string, error) {
			return "start_command", errors.New("Start command error")
		},
	}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
		getExecCommand: func(ExecOptions, *cli.App) error {
			return nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectForceRecreate(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "Start command error")
}

func TestCallStartProjectForceRecreateCase4(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "start_command up",
	}

	dlg := &testStartDialog{
		setMainContaner: func() (int, string, error) {
			return 0, "container_name", nil
		},
		startCommand: func() (string, error) {
			return "start_command", errors.New("Start command error")
		},
	}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
		getExecCommand: func(e ExecOptions, c *cli.App) error {
			assert.Equal(t, e.GetCommand(), "start_command")
			assert.Equal(t, e.GetArgs(), []string{"up", "--force-recreate", "container_name"})
			return nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{
		"container_name",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectForceRecreate(cfg, dlg, opt)

	assert.Nil(t, app.Action(ctx))
}

func TestCallStartProjectOrphansCase1(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "start_command up",
	}

	dlg := &testStartDialog{}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, errors.New("GetContainerList error")
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectOrphans(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "GetContainerList error")
}

func TestCallStartProjectOrphansCase2(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "start_command up",
	}

	dlg := &testStartDialog{
		setMainContaner: func() (int, string, error) {
			return 0, "", nil
		},
	}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectOrphans(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "Container name is empty. Set the container name")
}

func TestCallStartProjectOrphansCase3(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "",
	}

	dlg := &testStartDialog{
		setMainContaner: func() (int, string, error) {
			return 0, "container_name", nil
		},
		startCommand: func() (string, error) {
			return "start_command", errors.New("Start command error")
		},
	}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
		getExecCommand: func(ExecOptions, *cli.App) error {
			return nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectOrphans(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "Start command error")
}

func TestCallStartProjectOrphansCase4(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "start_command up",
	}

	dlg := &testStartDialog{
		setMainContaner: func() (int, string, error) {
			return 0, "container_name", nil
		},
		startCommand: func() (string, error) {
			return "start_command", errors.New("Start command error")
		},
	}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
		getExecCommand: func(e ExecOptions, c *cli.App) error {
			assert.Equal(t, e.GetCommand(), "start_command")
			assert.Equal(t, e.GetArgs(), []string{"up", "--remove-orphans", "container_name"})
			return nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{
		"container_name",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectOrphans(cfg, dlg, opt)

	assert.Nil(t, app.Action(ctx))
}

func TestCallStartProjectForceOrphansCase1(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "start_command up",
	}

	dlg := &testStartDialog{}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, errors.New("GetContainerList error")
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectForceOrphans(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "GetContainerList error")
}

func TestCallStartProjectForceOrphansCase2(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "start_command up",
	}

	dlg := &testStartDialog{
		setMainContaner: func() (int, string, error) {
			return 0, "", nil
		},
	}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectForceOrphans(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "Container name is empty. Set the container name")
}

func TestCallStartProjectForceOrphansCase3(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "",
	}

	dlg := &testStartDialog{
		setMainContaner: func() (int, string, error) {
			return 0, "container_name", nil
		},
		startCommand: func() (string, error) {
			return "start_command", errors.New("Start command error")
		},
	}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
		getExecCommand: func(ExecOptions, *cli.App) error {
			return nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectForceOrphans(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "Start command error")
}

func TestCallStartProjectForceOrphansCase4(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "start_command up",
	}

	dlg := &testStartDialog{
		setMainContaner: func() (int, string, error) {
			return 0, "container_name", nil
		},
		startCommand: func() (string, error) {
			return "start_command", errors.New("Start command error")
		},
	}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
		getExecCommand: func(e ExecOptions, c *cli.App) error {
			assert.Equal(t, e.GetCommand(), "start_command")
			assert.Equal(t, e.GetArgs(), []string{"up", "--force-recreate", "--remove-orphans", "container_name"})
			return nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{
		"container_name",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartProjectForceOrphans(cfg, dlg, opt)

	assert.Nil(t, app.Action(ctx))
}

func TestCallStartMainContainerCase1(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "start_command up",
	}

	dlg := &testStartDialog{}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, errors.New("GetContainerList error")
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartMainContainer(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "GetContainerList error")
}

func TestCallStartMainContainerCase2(t *testing.T) {
	cfg := &testStartConfig{
		getStartCommand: "start_command up",
	}

	dlg := &testStartDialog{
		setMainContaner: func() (int, string, error) {
			return 0, "", nil
		},
	}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartMainContainer(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "Container name is empty. Set the container name")
}

func TestCallStartMainContainerCase3(t *testing.T) {
	cfg := &testStartConfig{
		projectMainContainer: "container_name",
	}

	dlg := &testStartDialog{
		setMainContaner: func() (int, string, error) {
			return 0, "container_name", nil
		},
	}

	opt := &testStartOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{}, nil
		},
		getExecCommand: func(e ExecOptions, c *cli.App) error {
			assert.Equal(t, e.GetCommand(), "docker")
			assert.Equal(t, e.GetArgs(), []string{"start", "container_name"})
			return nil
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallStartMainContainer(cfg, dlg, opt)

	assert.Nil(t, app.Action(ctx))
}
