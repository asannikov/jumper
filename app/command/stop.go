package command

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2" // imports as package "cli"
)

type stopContainerOptions interface {
	GetInitFuntion() func(bool) string
	GetContainerList() ([]string, error)
	GetDockerStatus() bool
	GetStopContainers() func([]string) error
	GetExecCommand() func(ExecOptions, *cli.App) error
}

// CallStopAllContainersCommand stops all docker containers
func CallStopAllContainersCommand(options stopContainerOptions) *cli.Command {
	initf := options.GetInitFuntion()
	dockerStatus := options.GetDockerStatus()
	stopContainers := options.GetStopContainers()

	return &cli.Command{
		Name:    "stopallcontainers",
		Aliases: []string{"sac"},
		Usage:   "Stops all docker containers",
		Action: func(c *cli.Context) (err error) {
			initf(false)
			if dockerStatus {
				return stopContainers([]string{})
			}

			return errors.New("Docker is not running")
		},
	}
}

type callStopMainContainerCommandProjectConfig interface {
	GetProjectMainContainer() string
	SaveContainerNameToProjectConfig(string) error
}

type callStopMainContainerCommandDialog interface {
	SetMainContaner([]string) (int, string, error)
}

// CallStopMainContainerCommand stops main container
func CallStopMainContainerCommand(cfg callStopMainContainerCommandProjectConfig, d callStopMainContainerCommandDialog, options stopContainerOptions) *cli.Command {
	initf := options.GetInitFuntion()
	dockerStatus := options.GetDockerStatus()
	stopContainers := options.GetStopContainers()

	return &cli.Command{
		Name:    "stop:maincontainer",
		Aliases: []string{"smc"},
		Usage:   "Stops main docker container",
		Action: func(c *cli.Context) (err error) {
			if !dockerStatus {
				return errors.New("Docker is not running")
			}

			initf(true)

			var cl []string

			if cl, err = options.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			fmt.Printf("Searching for main container %s", cfg.GetProjectMainContainer())
			return stopContainers([]string{cfg.GetProjectMainContainer()})
		},
	}
}

// CallStopSelectedContainersCommand stops selected docker containers
func CallStopSelectedContainersCommand(options stopContainerOptions) *cli.Command {
	initf := options.GetInitFuntion()
	dockerStatus := options.GetDockerStatus()
	execCommand := options.GetExecCommand()

	return &cli.Command{
		Name:    "stop:containers",
		Aliases: []string{"scs"},
		Usage:   "Stops docker containers",
		Action: func(c *cli.Context) (err error) {
			initf(false)

			if !dockerStatus {
				return errors.New("Docker is not running")
			}

			eo := ExecOptions{
				command: "docker",
				args:    append([]string{"stop"}, c.Args().Slice()...),
				tty:     true,
				detach:  true,
			}

			return execCommand(eo, c.App)
		},
	}
}

// CallStopOneContainerCommand stops selected docker containers
// @todo
func CallStopOneContainerCommand(options stopContainerOptions) *cli.Command {
	initf := options.GetInitFuntion()
	dockerStatus := options.GetDockerStatus()
	execCommand := options.GetExecCommand()

	return &cli.Command{
		Name:    "stop:container",
		Aliases: []string{"stopc"},
		Usage:   "Stops selected docker containers",
		Action: func(c *cli.Context) (err error) {
			initf(false)

			if !dockerStatus {
				return errors.New("Docker is not running")
			}

			eo := ExecOptions{
				command: "docker",
				args:    append([]string{"stop"}, c.Args().Slice()...),
				tty:     true,
				detach:  true,
			}

			return execCommand(eo, c.App)
		},
	}
}
