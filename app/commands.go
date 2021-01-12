package app

import (
	"jumper/app/bash"    // github.com/asannikov/
	"jumper/app/command" // github.com/asannikov/
	"jumper/app/config"  // github.com/asannikov/
	"jumper/app/dialog"  // github.com/asannikov/
	"jumper/app/docker"  // github.com/asannikov/

	"github.com/urfave/cli/v2"
)

func commandList(c *config.Config, d *dialog.Dialog, initf func(bool) string) []*cli.Command {

	getCommandLocationF := bash.GetCommandLocation()

	dck := docker.GetDockerInstance()

	dockerDialog := getDockerStartDialog()
	dockerDialog.setDialog(&d)
	dockerDialog.setDocker(dck)
	dockerDialog.setDockerService(c.GetDockerCommand())

	dockerStatus := false

	if dockerAPIVersiongo, _ := dck.Stat(); dockerAPIVersiongo != "" {
		dockerStatus = true
		dck.InitClient()
	}

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
		command.SyncCommand("copyto", initf, dockerStatus, c, d, dockerDialog),
		command.SyncCommand("copyfrom", initf, dockerStatus, c, d, dockerDialog),

		// Xdebug
		command.XDebugCommand("xdebug:fpm:enable", initf, dockerStatus, c, d, dockerDialog),
		command.XDebugCommand("xdebug:fpm:disable", initf, dockerStatus, c, d, dockerDialog),
		command.XDebugCommand("xdebug:cli:enable", initf, dockerStatus, c, d, dockerDialog),
		command.XDebugCommand("xdebug:cli:disable", initf, dockerStatus, c, d, dockerDialog),

		// Shell
		command.ShellCommand(initf, c, d),

		// docker pull https://docs.docker.com/engine/api/sdk/examples/
	}
}
