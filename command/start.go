package command

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func defineStartCommand(cfg projectConfig, d dialog, containerlist []string) (err error) {
	if cfg.GetStartCommand() == "" {
		startCommand, err := d.StartCommand()

		if err != nil {
			return err
		}

		if startCommand == "" {
			return errors.New("Start command cannot be empty")
		}

		return cfg.SaveStartCommandToProjectConfig(startCommand)
	}

	return err
}

func runStartProject(c *cli.Context, cfg projectConfig, args []string) error {
	commandSlice := strings.Split(cfg.GetStartCommand(), " ")

	var binary = commandSlice[0]
	var initArgs = commandSlice[1:]

	extraInitArgs := c.Args().Slice()

	args = append(initArgs, args...)
	args = append(args, extraInitArgs...)

	log.Printf("Called: %s %s", binary, strings.Join(args, " "))

	cmd := exec.Command(binary, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// CallStartProjectBasic runs docker project
func CallStartProjectBasic(initf func(bool) string, cfg projectConfig, d dialog, clist containerlist) *cli.Command {
	cmd := cli.Command{
		Name:            "start",
		Aliases:         []string{"st"},
		Usage:           `runs defined command: {docker-compose -f docker-compose.yml up} [custom parameters]`,
		Description:     `It's possible to use any custom parameters coming after "up"`,
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf(true)

			var cl []string

			if cl, err = clist.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			if err = defineStartCommand(cfg, d, cl); err != nil {
				return err
			}

			return runStartProject(c, cfg, []string{})
		},
	}

	return &cmd
}

// CallStartProjectForceRecreate runs docker project
func CallStartProjectForceRecreate(initf func(bool) string, cfg projectConfig, d dialog, clist containerlist) *cli.Command {
	cmd := cli.Command{
		Name:    "start:force",
		Aliases: []string{"s:f"},
		Usage:   `runs defined command: {docker-compose -f docker-compose.yml up --force-recreat} [custom parameters]`,
		Description: `
		--force-recreate - Recreate containers even if their configuration and image haven't changed
		It's possible to use any custom parameters coming after "up"`,
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf(true)

			var cl []string

			if cl, err = clist.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			if err = defineStartCommand(cfg, d, cl); err != nil {
				return err
			}

			return runStartProject(c, cfg, []string{"--force-recreate"})
		},
	}

	return &cmd
}

// CallStartProjectOrphans runs docker project
func CallStartProjectOrphans(initf func(bool) string, cfg projectConfig, d dialog, clist containerlist) *cli.Command {
	cmd := cli.Command{
		Name:    "start:orphans",
		Aliases: []string{"s:o"},
		Usage:   `runs defined command: {docker-compose -f docker-compose.yml up --remove-orphans} [custom parameters]`,
		Description: `
		--remove-orphans - Remove containers for services not defined in the Compose file
		It's possible to use any custom parameters coming after "up"`,
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf(true)

			var cl []string

			if cl, err = clist.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			if err = defineStartCommand(cfg, d, cl); err != nil {
				return err
			}

			return runStartProject(c, cfg, []string{"--remove-orphans"})
		},
	}

	return &cmd
}

// CallStartProjectForceOrphans runs docker project
func CallStartProjectForceOrphans(initf func(bool) string, cfg projectConfig, d dialog, clist containerlist) *cli.Command {
	cmd := cli.Command{
		Name:    "start:force-orphans",
		Aliases: []string{"s:fo"},
		Usage:   `runs defined command: {docker-compose -f docker-compose.yml up --force-recreate --remove-orphans} [custom parameters]`,
		Description: `
		--force-recreate - Recreate containers even if their configuration and image haven't changed
		--remove-orphans - Remove containers for services not defined in the Compose file
		It's possible to use any custom parameters coming after "up"`,
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf(true)

			var cl []string

			if cl, err = clist.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			if err = defineStartCommand(cfg, d, cl); err != nil {
				return err
			}

			return runStartProject(c, cfg, []string{"--force-recreate", "--remove-orphans"})
		},
	}

	return &cmd
}

// CallStartMainContainer runs docker main container
func CallStartMainContainer(initf func(bool) string, cfg projectConfig, d dialog, clist containerlist) *cli.Command {
	cmd := cli.Command{
		Name:    "start:maincontainer",
		Aliases: []string{"startmc"},
		Usage:   `runs defined command: {docker start main_container}`,
		Action: func(c *cli.Context) (err error) {
			initf(true)

			var cl []string

			if cl, err = clist.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			args := []string{"start", cfg.GetProjectMainContainer()}
			fmt.Printf("\ncommand: %s\n\n", "docker "+strings.Join(args, " "))
			cmd := exec.Command("docker", args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		},
	}

	return &cmd
}

// CallRestartMainContainer restarts docker main container
func CallRestartMainContainer(initf func(bool) string, dockerStatus bool, cfg projectConfig, d dialog, clist containerlist) *cli.Command {
	cmd := cli.Command{
		Name:    "restart:maincontainer",
		Aliases: []string{"rmc"},
		Usage:   `restarts main container`,
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

			args := []string{"stop", cfg.GetProjectMainContainer()}
			fmt.Printf("\ncommand: %s\n\n", "docker "+strings.Join(args, " "))
			cmd := exec.Command("docker", args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}

			args = []string{"start", cfg.GetProjectMainContainer()}
			fmt.Printf("\ncommand: %s\n\n", "docker "+strings.Join(args, " "))
			cmd = exec.Command("docker", args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		},
	}

	return &cmd
}

// CallStartContainers runs docker custom container
func CallStartContainers(initf func(bool) string) *cli.Command {
	cmd := cli.Command{
		Name:    "start:containers",
		Aliases: []string{"startc"},
		Usage:   `runs defined command: {docker start} [container]`,
		Action: func(c *cli.Context) (err error) {
			initf(true)

			args := []string{"start"}
			args = append(args, c.Args().Slice()...)
			fmt.Printf("\ncommand: %s\n\n", "docker "+strings.Join(args, " "))
			cmd := exec.Command("docker", args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		},
	}

	return &cmd
}

// CallRestartContainers restart docker custom containers
func CallRestartContainers(initf func(bool) string, dockerStatus bool) *cli.Command {
	cmd := cli.Command{
		Name:    "restart:containers",
		Aliases: []string{"rc"},
		Usage:   `runs defined command: {docker start} [container]`,
		Action: func(c *cli.Context) (err error) {
			if !dockerStatus {
				return errors.New("Docker is not running")
			}

			initf(true)

			args := []string{"stop"}
			args = append(args, c.Args().Slice()...)
			fmt.Printf("\ncommand: %s\n\n", "docker "+strings.Join(args, " "))
			cmd := exec.Command("docker", args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return err
			}

			args = []string{"start"}
			args = append(args, c.Args().Slice()...)
			fmt.Printf("\ncommand: %s\n\n", "docker "+strings.Join(args, " "))
			cmd = exec.Command("docker", args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		},
	}

	return &cmd
}
