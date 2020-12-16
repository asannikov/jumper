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
		command.CallComposerCommand(initf, c, d, getContainerList()),
		command.CallComposerUpdateCommand(initf, c, d, getContainerList(), getCommandLocationF),
		command.CallComposerUpdateMemoryCommand(initf, c, d, getContainerList(), getCommandLocationF),

		//command.CallComposerInstallCommand(initf, c, d, getContainerList(), getCommandLocationF),
		//command.CallComposerInstallMemoryCommand(initf, c, d, getContainerList(), getCommandLocationF),

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
