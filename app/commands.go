package app

import (
	"os"

	"github.com/asannikov/jumper/app/config"
	"github.com/asannikov/jumper/app/dialog"

	"github.com/asannikov/jumper/app/command"

	"github.com/urfave/cli/v2"
)

type commandListDialog interface {
	DockerService() (string, error)
	StartDocker() (string, error)
	StartCommand() (string, error)
	SetMainContaner([]string) (int, string, error)
	DockerProjectPath(string) (string, error)
	CallAddProjectDialog(dialog.ProjectConfig) error
	AddProjectPath(string) (string, error)
	AddProjectName() (string, error)
	SelectProject([]string) (int, string, error)
	DockerShell() (int, string, error)
	DockerCliXdebugIniFilePath(string) (string, error)
	DockerFpmXdebugIniFilePath(string) (string, error)
	XDebugConfigLocation() (int, string, error)
}

type commandListOptions interface {
	GetInitFunction() func(bool) string
	GetCommandLocation() func(string, string) (string, error)
	GetStopContainers() func([]string) error
	GetExecCommand() func(command.ExecOptions, *cli.App) error
	GetDockerStatus() bool
	GetContainerList() ([]string, error)
	GetCopyTo(container string, sourcePath string, dstPath string) error
	RunNativeExec(eo command.ExecOptions, ca *cli.App) error
	DirExists(string) (bool, error)
	MkdirAll(string, os.FileMode) error
}

func commandList(c *config.Config, d commandListDialog, opt commandListOptions) []*cli.Command {
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
		command.ShellCommand(c, d, opt),

		// docker pull https://docs.docker.com/engine/api/sdk/examples/
		command.CallMagentoCommand(c, d, opt),
	}
}
