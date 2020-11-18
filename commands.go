package main

import (
	"mgt/command"

	"github.com/urfave/cli/v2"
)

func getCommandList(c *MgtConfig) []*cli.Command {

	phpContainer, _ := c.GetPhpContainer()

	return []*cli.Command{
		command.CallCliCommand(phpContainer),
		command.CallBashCommand(phpContainer),
		command.CallComposerCommand(phpContainer),
	}
}
