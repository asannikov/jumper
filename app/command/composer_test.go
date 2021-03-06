package command

import (
	"errors"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

type testComposerHandleBaseProjectConfig struct {
	mainContainer     string
	mainContainerUser string
	projectDockerPath string
}

func (tc *testComposerHandleBaseProjectConfig) GetProjectMainContainer() string {
	return tc.mainContainer
}

func (tc *testComposerHandleBaseProjectConfig) GetMainContainerUser() string {
	return tc.mainContainerUser
}

func (tc *testComposerHandleBaseProjectConfig) SaveContainerNameToProjectConfig(container string) error {
	return nil
}

func (tc *testComposerHandleBaseProjectConfig) SaveContainerUserToProjectConfig(user string) error {
	return nil
}

func (tc *testComposerHandleBaseProjectConfig) GetCommandInactveStatus(typecommand string) bool {
	return false
}

func (tc *testComposerHandleBaseProjectConfig) GetProjectDockerPath() string {
	return tc.projectDockerPath
}

type testComposerHandleBaseComposerDialog struct {
}

func (d *testComposerHandleBaseComposerDialog) SetMainContaner([]string) (int, string, error) {
	return 0, "", nil
}

func (d *testComposerHandleBaseComposerDialog) SetMainContanerUser() (string, error) {
	return "", nil
}

type testComposer struct {
	containerList []string
	locaton       func(string, string) (string, error)
	ctype         string
	command       string
}

func (tc *testComposer) GetCommandLocaton() func(string, string) (string, error) {
	return tc.locaton
}

func (tc *testComposer) GetCallType() string {
	return tc.ctype
}

func (tc *testComposer) GetComposerCommand() string {
	return tc.command
}

func (tc *testComposer) GetContainerList() []string {
	return tc.containerList
}

func TestParseCommand(t *testing.T) {

	shortcommand, calltype, dockercmd := parseCommand("composer:update:memory")
	assert.EqualValues(t, []string{"update:memory", "memory", "update"}, []string{shortcommand, calltype, dockercmd})

	shortcommand, calltype, dockercmd = parseCommand("composer:update")
	assert.EqualValues(t, []string{"update", "", "update"}, []string{shortcommand, calltype, dockercmd})

	shortcommand, calltype, dockercmd = parseCommand("composer:install:memory")
	assert.EqualValues(t, []string{"install:memory", "memory", "install"}, []string{shortcommand, calltype, dockercmd})

	shortcommand, calltype, dockercmd = parseCommand("composer:install")
	assert.EqualValues(t, []string{"install", "", "install"}, []string{shortcommand, calltype, dockercmd})

	shortcommand, calltype, dockercmd = parseCommand("composer")
	assert.EqualValues(t, []string{"composer", "", ""}, []string{shortcommand, calltype, dockercmd})

	shortcommand, calltype, dockercmd = parseCommand("composer:memory")
	assert.EqualValues(t, []string{"composer:memory", "memory", ""}, []string{shortcommand, calltype, dockercmd})
}

func TestComposerHandleCase1(t *testing.T) {
	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer: "",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	cmp := &testComposer{}

	a := &args{
		get:   "",
		slice: []string{},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{"container"},
	}

	_, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.EqualError(t, err, "Container name is empty. Set the container name")
}

func TestComposerHandleCase2(t *testing.T) {
	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer:     "containerName",
		mainContainerUser: "userName",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	cmp := &testComposer{}

	a := &args{
		get:   "",
		slice: []string{},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{"container"},
	}

	_, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.Nil(t, err)
}

func TestComposerHandleCase3(t *testing.T) {
	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer:     "containerName",
		mainContainerUser: "root",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	cmp := &testComposer{
		locaton: func(container string, service string) (string, error) {
			return "/path/to/" + service, nil
		},
		ctype:   "",
		command: "",
	}

	a := &args{
		get:   "",
		slice: []string{},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{"container"},
	}

	args, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.Nil(t, err)
	assert.Equal(t, []string{"exec", "-it", "containerName", "composer"}, args)
}

func TestComposerHandleCase4(t *testing.T) {
	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer:     "containerName",
		mainContainerUser: "root",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	cmp := &testComposer{
		locaton: func(container string, service string) (string, error) {
			return "/path/to/" + service, nil
		},
		ctype:   "",
		command: "",
	}

	a := &args{
		get: "",
		slice: []string{
			"update",
		},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{"container"},
	}

	args, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.Nil(t, err)
	assert.Equal(t, []string{"exec", "-it", "containerName", "composer", "update"}, args)
}

func TestComposerHandleCase5(t *testing.T) {
	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer:     "containerName",
		mainContainerUser: "root",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	cmp := &testComposer{
		locaton: func(container string, service string) (string, error) {
			return "/path/to/" + service, nil
		},
		ctype:   "",
		command: "",
	}

	a := &args{
		get: "m",
		slice: []string{
			"m",
			"update",
		},
		tail: []string{
			"update",
		},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{"container"},
	}

	args, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.Nil(t, err)
	assert.Equal(t, []string{"exec", "-i", "containerName", "/path/to/php", "-d", "memory_limit=-1", "/path/to/composer", "update"}, args)
}

func TestComposerHandleCase6(t *testing.T) {
	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer:     "containerName",
		mainContainerUser: "root",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	cmp := &testComposer{
		locaton: func(container string, service string) (string, error) {
			return "/path/to/" + service, nil
		},
		ctype:   "memory",
		command: "update",
	}

	a := &args{
		get: "",
		slice: []string{
			"--help",
		},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{"container"},
	}

	args, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.Nil(t, err)
	assert.Equal(t, []string{"exec", "-i", "containerName", "/path/to/php", "-d", "memory_limit=-1", "/path/to/composer", "update", "--help"}, args)
}

func TestComposerHandleCase7(t *testing.T) {
	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer:     "containerName",
		mainContainerUser: "userName",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	cmp := &testComposer{
		locaton: func(container string, service string) (string, error) {
			if service == "php" {
				return "", errors.New("Error on getting php path")
			}
			return "/path/to/" + service, nil
		},
		ctype:   "memory",
		command: "update",
	}

	a := &args{
		get: "",
		slice: []string{
			"--help",
		},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{"container"},
	}

	_, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.EqualError(t, err, "Error on getting php path")
}

func TestComposerHandleCase8(t *testing.T) {
	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer:     "containerName",
		mainContainerUser: "userName",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	cmp := &testComposer{
		locaton: func(container string, service string) (string, error) {
			if service == "composer" {
				return "", errors.New("Error on getting composer path")
			}
			return "/path/to/" + service, nil
		},
		ctype:   "memory",
		command: "update",
	}

	a := &args{
		get: "",
		slice: []string{
			"--help",
		},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{"container"},
	}

	_, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.EqualError(t, err, "Error on getting composer path")
}

func TestComposerHandleCase9(t *testing.T) {
	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer:     "containerName",
		mainContainerUser: "username",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	cmp := &testComposer{
		locaton: func(container string, service string) (string, error) {
			return "/path/to/" + service, nil
		},
		ctype:   "",
		command: "",
	}

	a := &args{
		get:   "",
		slice: []string{},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{"container"},
	}

	args, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.Nil(t, err)
	assert.Equal(t, []string{"exec", "-it", "-u", "username", "containerName", "composer"}, args)
}

func TestComposerHandleCase10(t *testing.T) {
	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer: "containerName",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	cmp := &testComposer{
		locaton: func(container string, service string) (string, error) {
			return "/path/to/" + service, nil
		},
		ctype:   "",
		command: "",
	}

	a := &args{
		get:   "",
		slice: []string{},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{"container"},
	}

	_, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.EqualError(t, err, "Container user name is empty. Set the user name")
}

type testComposerOption struct {
	getExecCommand     func(ExecOptions, *cli.App) error
	getInitFunction    func(bool) string
	getCommandLocation func(string, string) (string, error)
	getContainerList   func() ([]string, error)
}

func (x *testComposerOption) GetInitFunction() func(bool) string {
	return x.getInitFunction
}
func (x *testComposerOption) GetContainerList() ([]string, error) {
	return x.getContainerList()
}
func (x *testComposerOption) GetExecCommand() func(ExecOptions, *cli.App) error {
	return x.getExecCommand
}
func (x *testComposerOption) GetCommandLocation() func(string, string) (string, error) {
	return x.getCommandLocation
}

func TestCallComposerCommandCase1(t *testing.T) {
	opt := &testComposerOption{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getExecCommand: func(ExecOptions, *cli.App) error {
			return nil
		},
		getCommandLocation: func(string, string) (string, error) {
			return "local", nil
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, errors.New("composerHandle function error check")
		},
	}

	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer:     "containerName",
		mainContainerUser: "userName",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallComposerCommand("composer", cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "composerHandle function error check")
}

func TestCallComposerCommandCase2(t *testing.T) {
	opt := &testComposerOption{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getExecCommand: func(eo ExecOptions, a *cli.App) error {
			assert.Equal(t, eo.GetCommand(), "docker")
			assert.Equal(t, eo.GetArgs(), []string{"exec", "-it", "-u", "userName", "containerName", "composer"})
			return nil
		},
		getCommandLocation: func(string, string) (string, error) {
			return "local", nil
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
	}

	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer:     "containerName",
		mainContainerUser: "userName",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallComposerCommand("composer", cfg, dlg, opt)

	assert.Nil(t, app.Action(ctx))
}

func TestCallComposerCommandCase3(t *testing.T) {
	opt := &testComposerOption{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getExecCommand: func(eo ExecOptions, a *cli.App) error {
			assert.Equal(t, eo.GetCommand(), "docker")
			assert.Equal(t, eo.GetArgs(), []string{"exec", "-i", "-u", "userName", "containerName", "/usr/bin/php", "-d", "memory_limit=-1", "/usr/bin/composer"})
			return nil
		},
		getCommandLocation: func(v string, s string) (string, error) {
			if s == "php" {
				return "/usr/bin/php", nil
			}

			return "/usr/bin/composer", nil
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
	}

	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer:     "containerName",
		mainContainerUser: "userName",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallComposerCommand("composer:memory", cfg, dlg, opt)

	assert.Nil(t, app.Action(ctx))
}

func TestCallComposerCommandCase4(t *testing.T) {
	opt := &testComposerOption{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getExecCommand: func(eo ExecOptions, a *cli.App) error {
			assert.Equal(t, eo.GetCommand(), "docker")
			assert.Equal(t, eo.GetArgs(), []string{"exec", "-i", "-u", "userName", "containerName", "/usr/bin/php", "-d", "memory_limit=-1", "/usr/bin/composer", "install"})
			return nil
		},
		getCommandLocation: func(v string, s string) (string, error) {
			if s == "php" {
				return "/usr/bin/php", nil
			}

			return "/usr/bin/composer", nil
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
	}

	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer:     "containerName",
		mainContainerUser: "userName",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallComposerCommand("composer:install:memory", cfg, dlg, opt)

	assert.Nil(t, app.Action(ctx))
}

func TestCallComposerCommandCase5(t *testing.T) {
	opt := &testComposerOption{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getExecCommand: func(eo ExecOptions, a *cli.App) error {
			assert.Equal(t, eo.GetCommand(), "docker")
			assert.Equal(t, eo.GetArgs(), []string{"exec", "-i", "-u", "userName", "containerName", "/usr/bin/php", "-d", "memory_limit=-1", "/usr/bin/composer", "update"})
			return nil
		},
		getCommandLocation: func(v string, s string) (string, error) {
			if s == "php" {
				return "/usr/bin/php", nil
			}

			return "/usr/bin/composer", nil
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
	}

	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer:     "containerName",
		mainContainerUser: "userName",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallComposerCommand("composer:update:memory", cfg, dlg, opt)

	assert.Nil(t, app.Action(ctx))
}

func TestCallComposerCommandCase6(t *testing.T) {
	opt := &testComposerOption{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getExecCommand: func(eo ExecOptions, a *cli.App) error {
			assert.Equal(t, eo.GetCommand(), "docker")
			assert.Equal(t, eo.GetArgs(), []string{"exec", "-i", "-u", "userName", "containerName", "/usr/bin/php", "-d", "memory_limit=-1", "/usr/bin/composer", "update", "vvv"})
			return nil
		},
		getCommandLocation: func(v string, s string) (string, error) {
			if s == "php" {
				return "/usr/bin/php", nil
			}

			return "/usr/bin/composer", nil
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
	}

	cfg := &testComposerHandleBaseProjectConfig{
		mainContainer:     "containerName",
		mainContainerUser: "userName",
	}

	dlg := &testComposerHandleBaseComposerDialog{}

	set := &flag.FlagSet{}
	set.Parse([]string{"vvv"})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := CallComposerCommand("composer:update:memory", cfg, dlg, opt)

	assert.Nil(t, app.Action(ctx))
}
