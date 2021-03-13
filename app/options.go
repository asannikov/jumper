package app

import (
	"os"

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
	dirExists       func(string) (bool, error)
	mkdirAll        func(string, os.FileMode) error
	magentoBin      func(string, string) (bool, error)
	xdebugStatus    func(*cli.App, []string) (bool, error)
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

func (co *commandOptions) setDirExists(de func(string) (bool, error)) {
	co.dirExists = de
}

func (co *commandOptions) setMkdirAll(mka func(string, os.FileMode) error) {
	co.mkdirAll = mka
}

func (co *commandOptions) setMagentoBin(smb func(string, string) (bool, error)) {
	co.magentoBin = smb
}

func (co *commandOptions) setXdebugStatus(sxs func(*cli.App, []string) (bool, error)) {
	co.xdebugStatus = sxs
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

func (co *commandOptions) DirExists(path string) (bool, error) {
	return co.dirExists(path)
}

func (co *commandOptions) MkdirAll(path string, fileMode os.FileMode) error {
	return co.mkdirAll(path, fileMode)
}

func (co *commandOptions) CheckMagentoBin(containerName string, magentoBin string) (bool, error) {
	return co.magentoBin(containerName, magentoBin)
}

func (co *commandOptions) CheckXdebugStatus(ca *cli.App, args []string) (bool, error) {
	return co.xdebugStatus(ca, args)
}
