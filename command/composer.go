package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2" // imports as package "cli"
)

type composerConfig interface {
	GetProjectMainContainer() string
	SaveContainerNameToProjectConfig(string) error
}

type composerDialog interface {
	SetMainContaner([]string) (int, string, error)
}

type composer struct {
	usage         map[string]string
	aliases       map[string]string
	description   map[string]string
	initFunction  func()
	containerList []string
}

// CallComposerCommand generates composer commands
// https://medium.com/@ssttehrani/containers-from-scratch-with-golang-5276576f9909
// https://phase2.github.io/devtools/common-tasks/ssh-into-a-container/
func CallComposerCommand(composercommand string, initf func(), cfg projectConfig, d dialog, containerlist []string, getCommandLocation func(string, string) (string, error)) *cli.Command {
	commandstack := strings.Split(composercommand, ":")

	var calltype, dockercmd string
	if len(commandstack) == 2 {
		dockercmd = commandstack[1]
	}

	if len(commandstack) == 3 {
		calltype = commandstack[2]
	}

	shortcommand := strings.Join(commandstack[1:], ":")

	if shortcommand == "" {
		shortcommand = "composer"
	}

	cmp := composer{
		usage: map[string]string{
			"composer":        "Run composer [parameters]",
			"composer:memory": "Run composer [parameters]",
			"install":         "Run composer install",
			"install:memory":  "Run composer install",
			"update":          "Run composer update",
			"update:memory":   "Run composer update",
		},
		aliases: map[string]string{
			"composer":        "cmp",
			"composer:memory": "cmpm",
			"install":         "cmpi",
			"update":          "cmpu",
			"install:memory":  "cmpim",
			"update:memory":   "cmpum",
		},
		description: map[string]string{
			"composer":        "no text",
			"composer:memory": "no text",
			"install":         "Run composer install",
			"install:memory":  "Run composer install",
			"update":          "Run composer update",
			"update:memory":   "Run composer update",
		},
		containerList: containerlist,
	}

	return &cli.Command{
		Name:            composercommand,
		Aliases:         []string{cmp.aliases[shortcommand]},
		Usage:           cmp.usage[shortcommand],
		Description:     cmp.description[shortcommand],
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf()

			if err = defineProjectMainContainer(cfg, d, containerlist); err != nil {
				return err
			}

			var initArgs = []string{"exec"}

			composerArgs := []string{"-it", cfg.GetProjectMainContainer(), "composer", dockercmd}

			if calltype == "explicit" || c.Args().Get(0) == "m" {
				var phpLocation, composerLocation string

				if phpLocation, err = getCommandLocation(cfg.GetProjectMainContainer(), "php"); err != nil {
					return err
				}

				if composerLocation, err = getCommandLocation(cfg.GetProjectMainContainer(), "composer"); err != nil {
					return err
				}

				composerArgs = []string{"exec", "-i", cfg.GetProjectMainContainer(), phpLocation, "-d", "memory_limit=-1", composerLocation, dockercmd}
			}

			initArgs = append(initArgs, composerArgs...)

			extraInitArgs := c.Args().Slice()

			args := append(initArgs, extraInitArgs...)

			fmt.Printf("\n command: %s\n\n", "docker "+strings.Join(args, " "))

			cmd := exec.Command("docker", args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			return cmd.Run()
		},
	}
}
