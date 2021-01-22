package app

import (
	"log"

	"github.com/asannikov/jumper/app/bash"
	"github.com/asannikov/jumper/app/config"
	"github.com/asannikov/jumper/app/dialog"
	"github.com/asannikov/jumper/app/docker"

	"github.com/asannikov/jumper/app/command"

	"github.com/urfave/cli/v2"
)

func commandList(c *config.Config, d *dialog.Dialog, initf func(bool) string) []*cli.Command {

	b := bash.Bash{}

	getCommandLocationF := b.GetCommandLocation()

	dck := docker.GetDockerInstance()

	dockerDialog := getDockerStartDialog()
	dockerDialog.setDialog(d)
	dockerDialog.setDocker(dck)
	dockerDialog.setDockerService(c.GetDockerCommand())

	dockerStatus := false

	if dockerAPIVersiongo, _ := dck.Stat(); dockerAPIVersiongo != "" {
		dockerStatus = true
		dck.InitClient()
	}

	opt := &commandOptions{}
	opt.setInitFuntion(initf)
	opt.setCommandLocation(getCommandLocationF)
	opt.setDockerStatus(dockerStatus)
	opt.setStopContainers(dck.StopContainers())
	opt.setExecCommand(execCommand)
	opt.setDockerDialog(dockerDialog)
	opt.setCopyTo(func(container string, sourcePath string, dstPath string) error {
		return dck.CopyTo(container, sourcePath, dstPath)
	})
	opt.setNativeExec(func(container string, commands []string) (err error) {
		code, err := dck.Exec(container, commands)

		log.Println(code, err)
		return
	})

	return []*cli.Command{
		// cli commands
		command.CallCliCommand("cli", c, d, opt),
		command.CallCliCommand("sh", c, d, opt),
		command.CallCliCommand("clinotty", c, d, opt),
		command.CallCliCommand("cliroot", c, d, opt),
		command.CallCliCommand("clirootnotty", c, d, opt),

		// composer commands
		command.CallComposerCommand("composer", c, d, opt),
		command.CallComposerCommand("composer:memory", c, d, opt),
		command.CallComposerCommand("composer:install", c, d, opt),
		command.CallComposerCommand("composer:install:memory", c, d, opt),
		command.CallComposerCommand("composer:update", c, d, opt),
		command.CallComposerCommand("composer:update:memory", c, d, opt),

		// Docker start
		command.CallStartProjectBasic(c, d, opt),
		command.CallStartProjectForceRecreate(c, d, opt),
		command.CallStartProjectOrphans(c, d, opt),
		command.CallStartProjectForceOrphans(c, d, opt),
		command.CallStartMainContainer(c, d, opt),
		command.CallStartContainers(opt),

		// Docker restart
		command.CallRestartMainContainer(c, d, opt),
		command.CallRestartContainers(opt),

		// Stop all docker containers
		command.CallStopAllContainersCommand(opt),
		command.CallStopSelectedContainersCommand(opt),
		command.CallStopMainContainerCommand(c, d, opt),
		command.CallStopOneContainerCommand(opt),

		// Get Project Path
		command.GetProjectPath(d, opt),

		// Copyright
		command.CallCopyrightCommand(c, opt),

		// Sync Paths
		command.SyncCommand("copyto", c, d, opt),
		command.SyncCommand("copyfrom", c, d, opt),

		// Xdebug
		command.XDebugCommand("xdebug:fpm:enable", c, d, opt),
		command.XDebugCommand("xdebug:fpm:disable", c, d, opt),
		command.XDebugCommand("xdebug:cli:enable", c, d, opt),
		command.XDebugCommand("xdebug:cli:disable", c, d, opt),

		// Shell
		command.ShellCommand(initf, c, d),

		// docker pull https://docs.docker.com/engine/api/sdk/examples/
		command.CallMagentoCommand(c, d, opt),
	}
}
