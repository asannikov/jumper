package command

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2" // imports as package "cli"
	"os"
	"os/exec"
	"strings"
)

// CallStopAllContainersCommand stops all docker containers
func CallStopAllContainersCommand(initf func(bool) string, dockerStatus bool, stopFuncton func([]string) error) *cli.Command {
	return &cli.Command{
		Name:    "stopallcontainers",
		Aliases: []string{"sac"},
		Usage:   "Stops all docker containers",
		Action: func(c *cli.Context) (err error) {
			initf(false)
			if dockerStatus {
				return stopFuncton([]string{})
			}

			return errors.New("Docker is not running")
		},
	}
}

type callStopMainContainerCommandProjectConfig interface {
	GetProjectMainContainer() string
	SaveContainerNameToProjectConfig(string) error
}

// CallStopMainContainerCommand stops main container
func CallStopMainContainerCommand(initf func(bool) string, dockerStatus bool, stopFuncton func([]string) error, cfg callStopMainContainerCommandProjectConfig, d dialog, clist containerlist) *cli.Command {
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

			if cl, err = clist.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			fmt.Printf("Searching for main container %s", cfg.GetProjectMainContainer())
			return stopFuncton([]string{cfg.GetProjectMainContainer()})
		},
	}
}

// CallStopSelectedContainersCommand stops selected docker containers
func CallStopSelectedContainersCommand(initf func(bool) string, dockerStatus bool, stopFuncton func([]string) error) *cli.Command {
	return &cli.Command{
		Name:    "stop:containers",
		Aliases: []string{"scs"},
		Usage:   "Stops docker containers",
		Action: func(c *cli.Context) (err error) {
			initf(false)

			if !dockerStatus {
				return errors.New("Docker is not running")
			}

			args := []string{"stop"}

			args = append(args, c.Args().Slice()...)

			fmt.Printf("\ncommand: %s\n\n", "docker "+strings.Join(args, " "))

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
func CallStopOneContainerCommand(initf func(bool) string, dockerStatus bool, stopFuncton func([]string) error) *cli.Command {
	return &cli.Command{
		Name:    "stop:container",
		Aliases: []string{"stopc"},
		Usage:   "Stops selected docker containers",
		Action: func(c *cli.Context) (err error) {
			initf(false)

			if !dockerStatus {
				return errors.New("Docker is not running")
			}

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
