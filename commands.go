package main

import (
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

	return []*cli.Command{
		/*command.CallCliCommand(getPhpContainerName),
		command.CallBashCommand(getPhpContainerName),
		command.CallComposerCommand(getPhpContainerName),
		command.CallComposerUpdateCommand(getPhpContainerName),
		command.CallComposerUpdateMemoryCommand(getPhpContainerName),
		command.CallCopyFromContainer(getPhpContainerName),*/
		//command.CallComposerCommand1(getPhpContainerName),
		///Volumes/LS/LS/projects/docker/magento/m24/docker-compose.yml

		// Docker start
		command.CallStartProjectBasic(initf, c, d, getContainerList()),
		command.CallStartProjectForceRecreate(initf, c, d, getContainerList()),
		command.CallStartProjectOrphans(initf, c, d, getContainerList()),
		command.CallStartProjectForceOrphans(initf, c, d, getContainerList()),
	}
}
