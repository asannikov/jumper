package main

import (
	"mgt/command"

	"github.com/urfave/cli/v2"
)

func getCommandList(c *MgtConfig) []*cli.Command {

	getPhpContainerName := func() string {
		phpContainer, _ := c.GetPhpContainer()
		return phpContainer
	}

	return []*cli.Command{
		command.CallCliCommand(getPhpContainerName),
		command.CallBashCommand(getPhpContainerName),
		command.CallComposerCommand(getPhpContainerName),
		command.CallComposerUpdateCommand(getPhpContainerName),
		command.CallComposerUpdateMemoryCommand(getPhpContainerName),
		//command.CallComposerCommand1(getPhpContainerName),
	}
}
