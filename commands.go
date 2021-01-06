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
	DockerProjectPath(string) (string, error)
	DockerCliXdebugIniFilePath(string) (string, error)
	DockerFpmXdebugIniFilePath(string) (string, error)
	XDebugConfigLocation() (int, string, error)
}

func getCommandList(c *config.Config, d dialogCommand, initf func(bool) string) []*cli.Command {

	getCommandLocationF := bash.GetCommandLocation()

	dck := docker.GetDockerInstance()

	cl := getDockerStartDialog()
	cl.setDialog(d)
	cl.setDocker(dck)
	cl.setDockerService(c.GetDockerCommand())

	dockerStatus := false
	if dockerAPIVersiongo, _ := dck.Stat(); dockerAPIVersiongo != "" {
		dockerStatus = true
		dck.InitClient()
	}

	return []*cli.Command{
		// cli commands
		command.CallCliCommand("cli", initf, c, d, cl),
		command.CallCliCommand("sh", initf, c, d, cl),
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
		command.CallRestartMainContainer(initf, dockerStatus, c, d, cl),
		command.CallRestartContainers(initf, dockerStatus),

		// Stop all docker containers
		command.CallStopAllContainersCommand(initf, dockerStatus, dck.StopContainers()),
		command.CallStopSelectedContainersCommand(initf, dockerStatus, dck.StopContainers()),
		command.CallStopMainContainerCommand(initf, dockerStatus, dck.StopContainers(), c, d, cl),
		command.CallStopOneContainerCommand(initf, dockerStatus, dck.StopContainers()),

		// Get Project Path
		command.GetProjectPath(initf, d),

		// Copyright
		command.CallCopyrightCommand(initf, c, d),

		// Sync Paths
		command.SyncCommand("copyto", initf, dockerStatus, c, d, cl),
		command.SyncCommand("copyfrom", initf, dockerStatus, c, d, cl),

		// Xdebug
		command.XDebugCommand("xdebug:fpm:enable", initf, dockerStatus, c, d, cl),
		command.XDebugCommand("xdebug:fpm:disable", initf, dockerStatus, c, d, cl),
		command.XDebugCommand("xdebug:cli:enable", initf, dockerStatus, c, d, cl),
		command.XDebugCommand("xdebug:cli:disable", initf, dockerStatus, c, d, cl),

		// docker pull https://docs.docker.com/engine/api/sdk/examples/
	}
}
