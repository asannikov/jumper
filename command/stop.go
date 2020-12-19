package command

import (
	"fmt"
	"github.com/urfave/cli/v2" // imports as package "cli"
)

// CallStopAllContainersCommand stops all docker containers
func CallStopAllContainersCommand(stopFuncton func([]string) error) *cli.Command {
	return &cli.Command{
		Name:    "stopallcontainers",
		Aliases: []string{"sac"},
		Usage:   "Stops all docker containers",
		Action: func(c *cli.Context) (err error) {
			return stopFuncton([]string{})
		},
	}
}

// CallStopMainContainerCommand stops main container
func CallStopMainContainerCommand(stopFuncton func([]string) error, initf func(), cfg projectConfig, d dialog, containerList []string) *cli.Command {
	return &cli.Command{
		Name:    "stop:maincontainer",
		Aliases: []string{"smc"},
		Usage:   "Stops main docker container",
		Action: func(c *cli.Context) (err error) {
			initf()

			if err = defineProjectMainContainer(cfg, d, containerList); err != nil {
				return err
			}

			fmt.Printf("Searching for main container %s", cfg.GetProjectMainContainer())
			return stopFuncton([]string{cfg.GetProjectMainContainer()})
		},
	}
}

// CallStopSelectedContainersCommand stops selected docker containers
// @todo
func CallStopSelectedContainersCommand(stopFuncton func([]string) error, containerList []string) *cli.Command {
	return &cli.Command{
		Name:    "stop:containers",
		Aliases: []string{"scs"},
		Usage:   "Stops selected docker containers",
		Action: func(c *cli.Context) (err error) {
			return stopFuncton([]string{})
		},
	}
}

// CallStopOneContainerCommand stops selected docker containers
// @todo
func CallStopOneContainerCommand(stopFuncton func([]string) error, containerList []string) *cli.Command {
	return &cli.Command{
		Name:    "stop:container",
		Aliases: []string{"sc"},
		Usage:   "Stops selected docker containers",
		Action: func(c *cli.Context) (err error) {
			return stopFuncton([]string{})
		},
	}
}
