package command

import (
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

func TestCallStopAllContainersCommand(t *testing.T) {
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
