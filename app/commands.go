package app

import (
	"github.com/asannikov/jumper/app/bash"
	"github.com/asannikov/jumper/app/config"
	"github.com/asannikov/jumper/app/dialog"
	"github.com/asannikov/jumper/app/docker"

	"github.com/asannikov/jumper/app/command"

	"github.com/urfave/cli/v2"
)

func commandList(c *config.Config, d *dialog.Dialog, initf func(bool) string) []*cli.Command {

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
	opt.setExecCommand(execCommand)
	opt.setDockerDialog(dockerDialog)

	return []*cli.Command{
		// cli commands
		command.CallCliCommand("cli", initf, c, d, dockerDialog),
		command.CallCliCommand("sh", initf, c, d, dockerDialog),
		command.CallCliCommand("clinotty", initf, c, d, dockerDialog),
		command.CallCliCommand("cliroot", initf, c, d, dockerDialog),
		command.CallCliCommand("clirootnotty", initf, c, d, dockerDialog),

		// composer commands
		command.CallComposerCommand("composer", initf, c, d, dockerDialog, getCommandLocationF),
		command.CallComposerCommand("composer:memory", initf, c, d, dockerDialog, getCommandLocationF),
		command.CallComposerCommand("composer:install", initf, c, d, dockerDialog, getCommandLocationF),
		command.CallComposerCommand("composer:install:memory", initf, c, d, dockerDialog, getCommandLocationF),
		command.CallComposerCommand("composer:update", initf, c, d, dockerDialog, getCommandLocationF),
		command.CallComposerCommand("composer:update:memory", initf, c, d, dockerDialog, getCommandLocationF),

		// Docker start
		command.CallStartProjectBasic(initf, c, d, dockerDialog),
		command.CallStartProjectForceRecreate(initf, c, d, dockerDialog),
		command.CallStartProjectOrphans(initf, c, d, dockerDialog),
		command.CallStartProjectForceOrphans(initf, c, d, dockerDialog),
		command.CallStartMainContainer(initf, c, d, dockerDialog),
		command.CallStartContainers(initf),

		// Docker restart
		command.CallRestartMainContainer(initf, dockerStatus, c, d, dockerDialog),
		command.CallRestartContainers(initf, dockerStatus),

		// Stop all docker containers
		command.CallStopAllContainersCommand(initf, dockerStatus, dck.StopContainers()),
		command.CallStopSelectedContainersCommand(initf, dockerStatus, dck.StopContainers()),
		command.CallStopMainContainerCommand(initf, dockerStatus, dck.StopContainers(), c, d, dockerDialog),
		command.CallStopOneContainerCommand(initf, dockerStatus, dck.StopContainers()),

		// Get Project Path
		command.GetProjectPath(initf, d),

		// Copyright
		command.CallCopyrightCommand(initf, c, d),

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
