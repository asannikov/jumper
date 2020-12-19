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
	}
}
