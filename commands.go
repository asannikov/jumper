package main

import (
	"jumper/bash"
	"jumper/command"
	"jumper/config"
	"jumper/container"
	"jumper/dialog"

	"github.com/urfave/cli/v2"
)

func getCommandList(c *config.Config, d *dialog.Dialog, initf func()) []*cli.Command {

	getCommandLocationF := bash.GetCommandLocation()

	return []*cli.Command{
		// cli commands
		command.CallCliCommand("cli", initf, c, d, getContainerList()),
		command.CallCliCommand("bash", initf, c, d, getContainerList()),
		command.CallCliCommand("clinotty", initf, c, d, getContainerList()),
		command.CallCliCommand("cliroot", initf, c, d, getContainerList()),
		command.CallCliCommand("clirootnotty", initf, c, d, getContainerList()),

		// composer commands
		command.CallComposerCommand("composer", initf, c, d, getContainerList(), getCommandLocationF),
		command.CallComposerCommand("composer:memory", initf, c, d, getContainerList(), getCommandLocationF),
		command.CallComposerCommand("composer:install", initf, c, d, getContainerList(), getCommandLocationF),
		command.CallComposerCommand("composer:install:memory", initf, c, d, getContainerList(), getCommandLocationF),
		command.CallComposerCommand("composer:update", initf, c, d, getContainerList(), getCommandLocationF),
		command.CallComposerCommand("composer:update:memory", initf, c, d, getContainerList(), getCommandLocationF),

		/*command.CallCopyFromContainer(getPhpContainerName),*/

		// Docker start
		command.CallStartProjectBasic(initf, c, d, getContainerList()),
		command.CallStartProjectForceRecreate(initf, c, d, getContainerList()),
		command.CallStartProjectOrphans(initf, c, d, getContainerList()),
		command.CallStartProjectForceOrphans(initf, c, d, getContainerList()),
		command.CallStartMainContainer(initf, c, d, getContainerList()),
		command.CallStartContainers(initf),

		// Docker restart
		command.CallRestartMainContainer(initf, c, d, getContainerList()),
		command.CallRestartContainers(initf),
		
		// Stop all docker containers
		command.CallStopAllContainersCommand(container.StopContainers()),
		command.CallStopSelectedContainersCommand(container.StopContainers(), getContainerList()),
		command.CallStopMainContainerCommand(container.StopContainers(), initf, c, d, getContainerList()),
		command.CallStopOneContainerCommand(container.StopContainers(), getContainerList()),
	}
}
