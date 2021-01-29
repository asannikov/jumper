package app

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/asannikov/jumper/app/bash"
	"github.com/asannikov/jumper/app/config"
	"github.com/asannikov/jumper/app/dialog"
	"github.com/asannikov/jumper/app/docker"
	"github.com/docker/docker/api/types"

	"github.com/asannikov/jumper/app/command"

	"github.com/urfave/cli/v2"
)

type commandListDialog interface {
	DockerService() (string, error)
	StartDocker() (string, error)
	StartCommand() (string, error)
	SetMainContaner([]string) (int, string, error)
	DockerProjectPath(string) (string, error)
	CallAddProjectDialog(dialog.ProjectConfig) error
	AddProjectPath(string) (string, error)
	AddProjectName() (string, error)
	SelectProject([]string) (int, string, error)
	DockerShell() (int, string, error)
	DockerCliXdebugIniFilePath(string) (string, error)
	DockerFpmXdebugIniFilePath(string) (string, error)
	XDebugConfigLocation() (int, string, error)
}

func commandList(c *config.Config, d commandListDialog, initf func(bool) string) []*cli.Command {

	b := bash.Bash{}

	getCommandLocationF := b.GetCommandLocation()

	dck := docker.GetDockerInstance()

	dockerDialog := getDockerStartDialog()
	dockerDialog.setDialog(d)
	dockerDialog.setDocker(dck)
	dockerDialog.setDockerService(c.GetDockerCommand())

	dockerStatus := false

	if dockerAPIVersiongo, _ := dck.Stat(); dockerAPIVersiongo != "" {
		dockerStatus = true
		dck.InitClient()
	}

	opt := &commandOptions{}
	opt.setInitFuntion(initf)
	opt.setCommandLocation(getCommandLocationF)
	opt.setDockerStatus(dockerStatus)
	opt.setStopContainers(dck.StopContainers())
	opt.setExecCommand(func(eo command.ExecOptions, app *cli.App) error {

		cmd := exec.Command(eo.GetCommand(), eo.GetArgs()...)

		cmd.Stdin = app.Reader
		cmd.Stdout = app.Writer
		cmd.Stderr = app.ErrWriter

		fmt.Printf("\ncommand: %s\n\n", eo.GetCommand()+" "+strings.Join(eo.GetArgs(), " "))

		return cmd.Run()
	})
	opt.setDockerDialog(dockerDialog)
	opt.setCopyTo(func(container string, sourcePath string, dstPath string) error {
		return dck.CopyTo(container, sourcePath, dstPath)
	})
	opt.setNativeExec(func(eo command.ExecOptions, app *cli.App) (err error) {
		ic := &ioCli{
			reader:    app.Reader,
			writer:    app.Writer,
			errWriter: app.ErrWriter,
		}

		cnf := types.ExecConfig{
			AttachStderr: true,
			AttachStdin:  true,
			AttachStdout: true,
			User:         eo.GetUser(),
			Tty:          eo.GetTty(),
			Cmd:          append([]string{eo.GetCommand()}, eo.GetArgs()...),
			WorkingDir:   eo.GetWorkingDir(),
		}

		status, err := dck.Exec(c.GetProjectMainContainer(), &cnf, ic)

		if err != nil {
			return err
		}

		if status > 0 {
			return errors.New("Error is occurred on exec function")
		}

		return nil
	})

	return []*cli.Command{
		// cli commands
		command.CallCliCommand("cli", c, d, opt),
		command.CallCliCommand("sh", c, d, opt),
		command.CallCliCommand("clinotty", c, d, opt),
		command.CallCliCommand("cliroot", c, d, opt),
		command.CallCliCommand("clirootnotty", c, d, opt),

		// composer commands
		command.CallComposerCommand("composer", c, d, opt),
		command.CallComposerCommand("composer:memory", c, d, opt),
		command.CallComposerCommand("composer:install", c, d, opt),
		command.CallComposerCommand("composer:install:memory", c, d, opt),
		command.CallComposerCommand("composer:update", c, d, opt),
		command.CallComposerCommand("composer:update:memory", c, d, opt),

		// Docker start
		command.CallStartProjectBasic(c, d, opt),
		command.CallStartProjectForceRecreate(c, d, opt),
		command.CallStartProjectOrphans(c, d, opt),
		command.CallStartProjectForceOrphans(c, d, opt),
		command.CallStartMainContainer(c, d, opt),
		command.CallStartContainers(opt),

		// Docker restart
		command.CallRestartMainContainer(c, d, opt),
		command.CallRestartContainers(opt),

		// Stop all docker containers
		command.CallStopAllContainersCommand(opt),
		command.CallStopSelectedContainersCommand(opt),
		command.CallStopMainContainerCommand(c, d, opt),
		command.CallStopOneContainerCommand(opt),

		// Get Project Path
		command.GetProjectPath(d, opt),

		// Copyright
		command.CallCopyrightCommand(c, opt),

		// Sync Paths
		command.SyncCommand("copyto", c, d, opt),
		command.SyncCommand("copyfrom", c, d, opt),

		// Xdebug
		command.XDebugCommand("xdebug:fpm:enable", c, d, opt),
		command.XDebugCommand("xdebug:fpm:disable", c, d, opt),
		command.XDebugCommand("xdebug:cli:enable", c, d, opt),
		command.XDebugCommand("xdebug:cli:disable", c, d, opt),

		// Shell
		command.ShellCommand(initf, c, d),

		// docker pull https://docs.docker.com/engine/api/sdk/examples/
		command.CallMagentoCommand(c, d, opt),
	}
}
