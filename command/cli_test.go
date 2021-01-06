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

type testCli struct {
	args    map[string][]string
	command map[string]string
}

func (tc *testCli) GetCommand(cmd string) string {
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

type testCliHandleBaseProjectConfig struct {
	mainContainer string
}

func (tc *testCliHandleBaseProjectConfig) GetProjectMainContainer() string {
	return tc.mainContainer
}

func (tc *testCliHandleBaseProjectConfig) SaveContainerNameToProjectConfig(container string) error {
	return nil
}

func (tc *testCliHandleBaseProjectConfig) GetShell() string {
	return ""
}

type testCliHandleBaseComposerDialog struct {
}

func (d *testCliHandleBaseComposerDialog) SetMainContaner([]string) (int, string, error) {
	return 0, "", nil
}

func TestCliHandleCase1(t *testing.T) {
	cfg := &testCliHandleBaseProjectConfig{
		mainContainer: "",
	}

	dlg := &testCliHandleBaseComposerDialog{}

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
	cfg := &testCliHandleBaseProjectConfig{
		mainContainer: "containerName",
	}

	dlg := &testCliHandleBaseComposerDialog{}

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
	cfg := &testCliHandleBaseProjectConfig{
		mainContainer: "containerName",
	}

	dlg := &testCliHandleBaseComposerDialog{}

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
