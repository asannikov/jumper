package app

import "github.com/urfave/cli/v2"

type commandOptions struct {
	initf           func(bool) string
	commandLocation func(string, string) (string, error)
	stopContainers  func([]string) error
	execCommand     func(string, []string, *cli.App) error
	dockerStatus    bool
	dockerDialog    *dockerStartDialog
}

func (co *commandOptions) setInitFuntion(f func(bool) string) {
	co.initf = f
}

func (co *commandOptions) setCommandLocation(cl func(string, string) (string, error)) {
	co.commandLocation = cl
}

func (co *commandOptions) setStopContainers(sc func([]string) error) {
	co.stopContainers = sc
}

func (co *commandOptions) setExecCommand(ec func(string, []string, *cli.App) error) {
	co.execCommand = ec
}

func (co *commandOptions) setDockerStatus(status bool) {
	co.dockerStatus = status
}

func (co *commandOptions) setDockerDialog(dd *dockerStartDialog) {
	co.dockerDialog = dd
}

func (co *commandOptions) GetInitFuntion() func(bool) string {
	return co.initf
}

func (co *commandOptions) GetCommandLocation() func(string, string) (string, error) {
	return co.commandLocation
}

func (co *commandOptions) GetStopContainers() func([]string) error {
	return co.stopContainers
}

func (co *commandOptions) GetExecCommand() func(string, []string, *cli.App) error {
	return co.execCommand
}

func (co *commandOptions) GetDockerStatus() bool {
	return co.dockerStatus
}

func (co *commandOptions) GetContainerList() ([]string, error) {
	return co.dockerDialog.GetContainerList()
}
