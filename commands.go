package main

import (
	"jumper/bash"
	"jumper/command"
	"jumper/config"
	"jumper/docker"

	"github.com/urfave/cli/v2"
)

type dialogCommand interface {
	SetMainContaner([]string) (int, string, error)
	StartCommand() (string, error)
	StartDocker() (string, error)
	DockerService() (string, error)
}

func getCommandList(c *config.Config, d dialogCommand, initf func(bool)) []*cli.Command {

	getCommandLocationF := bash.GetCommandLocation()

	dck := docker.GetDockerInstance()

	cl := getDockerStartDialog()
	cl.setDialog(d)
	cl.setDocker(dck)
	cl.setDockerService(c.GetDockerCommand())

	return []*cli.Command{
		// cli commands
		command.CallCliCommand("cli", initf, c, d, cl),
		command.CallCliCommand("bash", initf, c, d, cl),
		command.CallCliCommand("clinotty", initf, c, d, cl),
		command.CallCliCommand("cliroot", initf, c, d, cl),
		command.CallCliCommand("clirootnotty", initf, c, d, cl),

		// composer commands
		command.CallComposerCommand("composer", initf, c, d, cl, getCommandLocationF),
		command.CallComposerCommand("composer:memory", initf, c, d, cl, getCommandLocationF),
		command.CallComposerCommand("composer:install", initf, c, d, cl, getCommandLocationF),
		command.CallComposerCommand("composer:install:memory", initf, c, d, cl, getCommandLocationF),
		command.CallComposerCommand("composer:update", initf, c, d, cl, getCommandLocationF),
		command.CallComposerCommand("composer:update:memory", initf, c, d, cl, getCommandLocationF),

		// Docker start
		command.CallStartProjectBasic(initf, c, d, cl),
		command.CallStartProjectForceRecreate(initf, c, d, cl),
		command.CallStartProjectOrphans(initf, c, d, cl),
		command.CallStartProjectForceOrphans(initf, c, d, cl),
		command.CallStartMainContainer(initf, c, d, cl),
		command.CallStartContainers(initf),

		// Docker restart
		command.CallRestartMainContainer(initf, c, d, cl),
		command.CallRestartContainers(initf),

		// Stop all docker containers
		command.CallStopAllContainersCommand(dck.StopContainers()),
		command.CallStopSelectedContainersCommand(dck.StopContainers()),
		command.CallStopMainContainerCommand(dck.StopContainers(), initf, c, d, cl),
		command.CallStopOneContainerCommand(dck.StopContainers()),

		// Get Project Path
		command.GetProjectPath(initf, c, d),

		// Copyright
		command.CallCopyrightCommand(initf, c, d),

		// docker pull https://docs.docker.com/engine/api/sdk/examples/
	}
}
