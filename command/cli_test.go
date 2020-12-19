package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

type testCliConfig struct {
	mainContainer string
}

func (tc *testCliConfig) GetProjectMainContainer() string {
	return tc.mainContainer
}

func (tc *testCliConfig) SaveContainerNameToProjectConfig(container string) error {
	return nil
}

func (tc *testCliConfig) GetStartCommand() string {
	return ""
}

func (tc *testCliConfig) SaveStartCommandToProjectConfig(c string) error {
	return nil
}

type testCliDialog struct{}

func (d *testCliDialog) SetMainContaner([]string) (int, string, error) {
	return 0, "", nil
}

func (d *testCliDialog) SetStartCommand() (string, error) {
	return "", nil
}

type testCli struct {
	containerList []string
	args          map[string][]string
	command       map[string]string
}

func (tc *testCli) GetCommand(cmd string) string {
	return tc.command[cmd]
}

func (tc *testCli) GetArgs() map[string][]string {
	return tc.args
}

func (tc *testCli) GetContainerList() []string {
	return tc.containerList
}

func TestCliHandleCase1(t *testing.T) {
	cfg := &testCliConfig{
		mainContainer: "",
	}

	dlg := &testCliDialog{}

	cli := &testCli{}

	a := &args{
		get:   "",
		slice: []string{},
	}

	_, err := cliCommandHandle("cli", cfg, dlg, cli, a)

	assert.EqualError(t, err, "Container name is empty. Set the container name")
}

func TestCliHandleCase2(t *testing.T) {
	cfg := &testCliConfig{
		mainContainer: "containerName",
	}

	dlg := &testCliDialog{}

	cli := &testCli{}

	a := &args{
		get:   "",
		slice: []string{},
	}

	_, err := cliCommandHandle("cli", cfg, dlg, cli, a)

	assert.EqualError(t, err, "Please specify a CLI command (ex. ls)")
}

func TestCliHandleCase3(t *testing.T) {
	cfg := &testCliConfig{
		mainContainer: "containerName",
	}

	dlg := &testCliDialog{}

	cli := &testCli{
		args: map[string][]string{
			"cli": []string{"-it"},
		},
		command: map[string]string{
			"cli": "bash",
		},
	}

	a := &args{
		get:   "",
		slice: []string{},
	}

	args, err := cliCommandHandle("cli", cfg, dlg, cli, a)

	assert.Nil(t, err)
	assert.Equal(t, []string{"exec", "-it", "containerName", "bash"}, args)
}
