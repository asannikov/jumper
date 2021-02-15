package command

import (
	"errors"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

type testShellConfig struct {
	shell string
	save  error
}

func (s *testShellConfig) GetShell() string {
	return s.shell
}

func (s *testShellConfig) SaveShellCommand(v string) error {
	return s.save
}

type testShellDialog struct {
	shell string
	err   error
}

func (s *testShellDialog) DockerShell() (int, string, error) {
	return 0, s.shell, s.err
}

type testShellOptions struct {
	getInitFunction func(bool) string
}

func (x *testShellOptions) GetInitFunction() func(bool) string {
	return x.getInitFunction
}

func TestDefineShellTypeCase1(t *testing.T) {
	cfg := &testShellConfig{}
	d := &testShellDialog{
		shell: "bash",
		err:   nil,
	}

	assert.Nil(t, defineShellType(cfg, d))
}

func TestDefineShellTypeCase2(t *testing.T) {
	cfg := &testShellConfig{}
	d := &testShellDialog{
		shell: "",
		err:   errors.New("Something goes wrong. Shell was not set"),
	}

	assert.EqualError(t, defineShellType(cfg, d), "Something goes wrong. Shell was not set")
}

func TestDefineShellTypeCase3(t *testing.T) {
	cfg := &testShellConfig{
		save: errors.New("cannot save shell"),
	}
	d := &testShellDialog{
		shell: "bash",
		err:   nil,
	}

	assert.EqualError(t, defineShellType(cfg, d), "cannot save shell")
}

func TestDefineShellTypeCase4(t *testing.T) {
	cfg := &testShellConfig{
		save: errors.New("cannot save shell"),
	}
	d := &testShellDialog{
		shell: "",
		err:   nil,
	}

	assert.EqualError(t, defineShellType(cfg, d), "Something goes wrong. Shell was not set")
}

func TestDefineShellTypeCase5(t *testing.T) {
	cfg := &testShellConfig{}
	dlg := &testShellDialog{
		shell: "",
		err:   errors.New("DockerShell error"),
	}
	opt := &testShellOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := ShellCommand(cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "DockerShell error")
}

func TestDefineShellTypeCase6(t *testing.T) {
	cfg := &testShellConfig{}
	dlg := &testShellDialog{
		shell: "path",
		err:   nil,
	}
	opt := &testShellOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := ShellCommand(cfg, dlg, opt)

	assert.Nil(t, app.Action(ctx))
}
