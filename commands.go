package main

import (
	"mgt/bash"
	"mgt/command"
	"mgt/config"
	"mgt/dialog"

	"github.com/urfave/cli/v2"
)

func getCommandList(c *config.Config, d *dialog.Dialog, initf func()) []*cli.Command {

	/*getPhpContainerName := func() (string, error) {
		phpContainer, err := c.GetPhpContainer(func() error {
			return c.EvaluatePhpContainer(func() (int, string, error) {
				return d.SetPhpContaner(getContainerList())
			})
		})
		return phpContainer, err
	}*/

	getCommandLocationF := bash.GetCommandLocation()

	return []*cli.Command{
		// composer commands
		command.CallComposerCommand("composer", initf, c, d, getContainerList(), getCommandLocationF),
		command.CallComposerCommand("composer:memory", initf, c, d, getContainerList(), getCommandLocationF),
		command.CallComposerCommand("composer:install", initf, c, d, getContainerList(), getCommandLocationF),
		command.CallComposerCommand("composer:install:memory", initf, c, d, getContainerList(), getCommandLocationF),
		command.CallComposerCommand("composer:update", initf, c, d, getContainerList(), getCommandLocationF),
		command.CallComposerCommand("composer:update:memory", initf, c, d, getContainerList(), getCommandLocationF),

		/*command.CallCliCommand(getPhpContainerName),
		command.CallBashCommand(getPhpContainerName),

		command.CallComposerUpdateMemoryCommand(getPhpContainerName),
		command.CallCopyFromContainer(getPhpContainerName),*/
		//command.CallComposerCommand1(getPhpContainerName),

		// Docker start
		command.CallStartProjectBasic(initf, c, d, getContainerList()),
		command.CallStartProjectForceRecreate(initf, c, d, getContainerList()),
		command.CallStartProjectOrphans(initf, c, d, getContainerList()),
		command.CallStartProjectForceOrphans(initf, c, d, getContainerList()),
	}
}
