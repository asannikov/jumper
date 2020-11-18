package main

import (
	"mgt/command"

	"github.com/urfave/cli/v2"
)

func getCommandList(c *MgtConfig) []*cli.Command {
	return []*cli.Command{
		command.CallCliCommand(),
		command.CallBashCommand(c.GetPhpContainer),
		command.CallComposerCommand(),
	}
}
