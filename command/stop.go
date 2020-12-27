package command

import (
	"fmt"
	"github.com/urfave/cli/v2" // imports as package "cli"
	"os"
	"os/exec"
	"strings"
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
func CallStopMainContainerCommand(stopFuncton func([]string) error, initf func(bool), cfg projectConfig, d dialog, containerList []string) *cli.Command {
	return &cli.Command{
		Name:    "stop:maincontainer",
		Aliases: []string{"smc"},
		Usage:   "Stops main docker container",
		Action: func(c *cli.Context) (err error) {
			initf(true)

			if err = defineProjectMainContainer(cfg, d, containerList); err != nil {
				return err
			}

			fmt.Printf("Searching for main container %s", cfg.GetProjectMainContainer())
			return stopFuncton([]string{cfg.GetProjectMainContainer()})
		},
	}
}

// CallStopSelectedContainersCommand stops selected docker containers
func CallStopSelectedContainersCommand(stopFuncton func([]string) error, containerList []string) *cli.Command {
	return &cli.Command{
		Name:    "stop:containers",
		Aliases: []string{"scs"},
		Usage:   "Stops docker containers",
		Action: func(c *cli.Context) (err error) {
			args := []string{"stop"}

			args = append(args, c.Args().Slice()...)

			fmt.Printf("\n command: %s\n\n", "docker "+strings.Join(args, " "))

			cmd := exec.Command("docker", args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		},
	}
}

// CallStopOneContainerCommand stops selected docker containers
// @todo
func CallStopOneContainerCommand(stopFuncton func([]string) error, containerList []string) *cli.Command {
	return &cli.Command{
		Name:    "stop:container",
		Aliases: []string{"stopc"},
		Usage:   "Stops selected docker containers",
		Action: func(c *cli.Context) (err error) {
			args := []string{"stop"}

			args = append(args, c.Args().Slice()...)
			cmd := exec.Command("docker", args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		},
	}
}
