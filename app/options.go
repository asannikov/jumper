package app

import (
	"github.com/asannikov/jumper/app/command"
	"github.com/urfave/cli/v2"
)

type commandOptionsDockerDialog interface {
	GetContainerList() ([]string, error)
}

type commandOptions struct {
	initf           func(bool) string
	commandLocation func(string, string) (string, error)
	stopContainers  func([]string) error
	execCommand     func(command.ExecOptions, *cli.App) error
	copyTo          func(string, string, string) error
	dockerStatus    bool
	dockerDialog    commandOptionsDockerDialog
	nativeExec      func(command.ExecOptions, *cli.App) (err error)
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

func (co *commandOptions) setExecCommand(ec func(command.ExecOptions, *cli.App) error) {
	co.execCommand = ec
}

func (co *commandOptions) setDockerStatus(status bool) {
	co.dockerStatus = status
}

func (co *commandOptions) setDockerDialog(dd commandOptionsDockerDialog) {
	co.dockerDialog = dd
}

func (co *commandOptions) setCopyTo(ct func(string, string, string) error) {
	co.copyTo = ct
}

func (co *commandOptions) setNativeExec(ne func(command.ExecOptions, *cli.App) (err error)) {
	co.nativeExec = ne
}

func (co *commandOptions) GetInitFunction() func(bool) string {
	return co.initf
}

func (co *commandOptions) GetCommandLocation() func(string, string) (string, error) {
	return co.commandLocation
}

func (co *commandOptions) GetStopContainers() func([]string) error {
	return co.stopContainers
}

func (co *commandOptions) GetExecCommand() func(command.ExecOptions, *cli.App) error {
	return co.execCommand
}

func (co *commandOptions) GetDockerStatus() bool {
	return co.dockerStatus
}

func (co *commandOptions) GetContainerList() ([]string, error) {
	return co.dockerDialog.GetContainerList()
}

func (co *commandOptions) GetCopyTo(container string, sourcePath string, dstPath string) error {
	return co.copyTo(container, sourcePath, dstPath)
}

func (co *commandOptions) RunNativeExec(eo command.ExecOptions, ca *cli.App) error {
	return co.nativeExec(eo, ca)
}
