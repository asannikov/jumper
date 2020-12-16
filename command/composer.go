package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2" // imports as package "cli"
)

type compoaserConfig interface {
	GetProjectMainContainer() string
	SaveContainerNameToProjectConfig(string) error
}

type compoaserDialog interface {
	SetMainContaner([]string) (int, string, error)
}

// CallComposerCommand generates Composer command (docker-compose exec phpfpm composer)
func CallComposerCommand(initf func(), cfg startProjectConfig, d startDialog, containerlist []string) *cli.Command {

	cmd := cli.Command{
		Name:            "composer",
		Aliases:         []string{"cmp"},
		Usage:           "Run composer",
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf()

			if err = defineProjectMainContainer(cfg, d, containerlist); err != nil {
				return err
			}

			var initArgs = []string{"exec", "-it", cfg.GetProjectMainContainer(), "composer"}

			extraInitArgs := c.Args().Slice()

			args := append(initArgs, extraInitArgs...)

			fmt.Printf("\n command: %s\n\n", "docker "+strings.Join(args, " "))

			cmd := exec.Command("docker", args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			cmd.Run()

			return nil
		},
	}

	return &cmd
}

// CallComposerUpdateCommand generates Composer command (docker-compose exec phpfpm composer)
func CallComposerUpdateCommand(initf func(), cfg startProjectConfig, d startDialog, containerlist []string, getCommandLocation func(string, string) (string, error)) *cli.Command {

	cmd := cli.Command{
		Name:            "composer:update",
		Aliases:         []string{"cmpu"},
		Usage:           "Run composer update",
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf()

			if err = defineProjectMainContainer(cfg, d, containerlist); err != nil {
				return err
			}

			if c.Args().Get(0) == "m" {
				return explictComposerUpdate(cfg.GetProjectMainContainer(), c.Args().Tail(), getCommandLocation)
			}

			var initArgs = []string{"exec", "-it", cfg.GetProjectMainContainer(), "composer", "update"}

			extraInitArgs := c.Args().Slice()

			args := append(initArgs, extraInitArgs...)
			// https://medium.com/@ssttehrani/containers-from-scratch-with-golang-5276576f9909
			// https://phase2.github.io/devtools/common-tasks/ssh-into-a-container/

			fmt.Printf("\n command: %s\n\n", "docker "+strings.Join(args, " "))

			cmd := exec.Command("docker", args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			cmd.Run()

			return nil
		},
	}

	return &cmd
}

// CallComposerUpdateMemoryCommand generates Composer command (docker-compose exec phpfpm php -d memory_limit composer update)
func CallComposerUpdateMemoryCommand(initf func(), cfg startProjectConfig, d startDialog, containerlist []string, getCommandLocation func(string, string) (string, error)) *cli.Command {

	cmd := cli.Command{
		Name:            "composer:update:memory",
		Aliases:         []string{"cmpum"},
		Usage:           "Run php -d memory_limit=-1 composer update",
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf()

			if err = defineProjectMainContainer(cfg, d, containerlist); err != nil {
				return err
			}

			return explictComposerUpdate(cfg.GetProjectMainContainer(), c.Args().Slice(), getCommandLocation)
		},
	}

	return &cmd
}

func explictComposerUpdate(phpContainer string, extraInitArgs []string, getCommandLocation func(string, string) (string, error)) (err error) {
	var phpLocation, composerLocation string

	if phpLocation, err = getCommandLocation(phpContainer, "php"); err != nil {
		return err
	}

	if composerLocation, err = getCommandLocation(phpContainer, "composer"); err != nil {
		return err
	}

	var initArgs = []string{"exec", "-i", phpContainer, phpLocation, "-d", "memory_limit=-1", composerLocation, "update"}

	args := append(initArgs, extraInitArgs...)

	fmt.Printf("\n command: %s\n\n", "docker "+strings.Join(args, " "))

	cmd := exec.Command("docker", args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
