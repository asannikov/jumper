package main

import (
	"mgt/command"
	"mgt/config"

	"github.com/urfave/cli/v2"
)

func getCommandList(c *config.MgtConfig, d *dialogInternal) []*cli.Command {

	getPhpContainerName := func() (string, error) {
		phpContainer, err := c.GetPhpContainer(func() error {
			return c.GenerateConfig(func() error {
				return c.EvaluatePhpContainer(func() (int, string, error) {
					return d.setPhpContaner()
				})
			})
		})
		return phpContainer, err
	}

	getProjectPath := func() (string, error) {
		projectPath, err := c.GetProjectPath(func() error {
			return c.GenerateConfig(func() error {
				return c.EvaluateProjectPath(func(projects []string) (int, string, error) {
					return d.selectProject(projects)
				}, func() (string, error) {
					return d.addProjectName()
				}, func(path string) (string, error) {
					return d.addProjectPath(path)
				})
			})
		})
		return projectPath, err
	}

	return []*cli.Command{
		command.CallCliCommand(getPhpContainerName),
		command.CallBashCommand(getPhpContainerName),
		command.CallComposerCommand(getPhpContainerName),
		command.CallComposerUpdateCommand(getPhpContainerName),
		command.CallComposerUpdateMemoryCommand(getPhpContainerName),
		command.CallCopyFromContainer(getPhpContainerName),
		//command.CallComposerCommand1(getPhpContainerName),
		///Volumes/LS/LS/projects/docker/magento/m24/docker-compose.yml
		command.CallStartProject(getProjectPath),
	}
}
