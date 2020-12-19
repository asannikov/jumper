package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2" // imports as package "cli"
)

type composer struct {
	usage         map[string]string
	aliases       map[string]string
	description   map[string]string
	initFunction  func()
	containerList []string
	locaton       func(string, string) (string, error)
	ctype         string
	command       string
}

func (cmp *composer) GetCommandLocaton() func(string, string) (string, error) {
	return cmp.locaton
}

func (cmp *composer) GetCallType() string {
	return cmp.ctype
}

func (cmp *composer) GetComposerCommand() string {
	return cmp.command
}

func (cmp *composer) GetContainerList() []string {
	return cmp.containerList
}

func parseCommand(composercommand string) (string, string, string) {
	commandstack := strings.Split(composercommand, ":")

	var calltype, dockercmd string
	if len(commandstack) > 1 {
		dockercmd = strings.Trim(commandstack[1], " ")
	}

	if len(commandstack) > 2 {
		calltype = commandstack[2]
	}

	index := strings.Join(commandstack[1:], ":")

	if index == "" {
		index = "composer"
	} else if index == "memory" {
		index = "composer:memory"
		dockercmd = ""
		calltype = "memory"
	}

	return index, calltype, dockercmd
}

// CallComposerCommand generates composer commands
// https://medium.com/@ssttehrani/containers-from-scratch-with-golang-5276576f9909
// https://phase2.github.io/devtools/common-tasks/ssh-into-a-container/
func CallComposerCommand(composercommand string, initf func(), cfg projectConfig, d dialog, containerlist []string, getCommandLocation func(string, string) (string, error)) *cli.Command {
	index, calltype, dockercmd := parseCommand(composercommand)

	cmp := &composer{
		usage: map[string]string{
			"composer":        "Runs composer: {docker exec -it phpContainer composer} [custom parameters]",
			"composer:memory": "Runs composer with no memory constraint: {docker exec -i phpContainer /usr/bin/php -d memory_limit=-1 /usr/local/bin/composer} [custom parameters]",
			"install":         "Runs composer install: {docker exec -it phpContainer composer install} [custom parameters]",
			"install:memory":  "Runs composer install with no memory constraint: {docker exec -i phpContainer /usr/bin/php -d memory_limit=-1 /usr/local/bin/composer install} [custom parameters]",
			"update":          "Runs composer update: {docker exec -it phpContainer composer update} [custom parameters]",
			"update:memory":   "Runs composer update with no memory constraint: {docker exec -i phpContainer /usr/bin/php -d memory_limit=-1 /usr/local/bin/composer update} [custom parameters]",
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
			"composer":        "phpContainer is taken from project config file",
			"composer:memory": "phpContainer is taken from project config file, php and composer commands will be found automatically",
			"install":         "phpContainer is taken from project config file",
			"install:memory":  "phpContainer is taken from project config file, php and composer commands will be found automatically",
			"update":          "phpContainer is taken from project config file",
			"update:memory":   "phpContainer is taken from project config file, php and composer commands will be found automatically",
		},
		containerList: containerlist,
		locaton:       getCommandLocation,
		ctype:         calltype,
		command:       dockercmd,
	}

	return &cli.Command{
		Name:            composercommand,
		Aliases:         []string{cmp.aliases[index]},
		Usage:           cmp.usage[index],
		Description:     cmp.description[index],
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf()

			var args []string

			if args, err = composerHandle(cfg, d, cmp, c.Args()); err != nil {
				return err
			}

			fmt.Printf("\n command: %s\n\n", "docker "+strings.Join(args, " "))

			cmd := exec.Command("docker", args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			return cmd.Run()
		},
	}
}

type composerInterface interface {
	GetCommandLocaton() func(string, string) (string, error)
	GetCallType() string
	GetComposerCommand() string
	GetContainerList() []string
}

func composerHandle(cfg projectConfig, d dialog, c composerInterface, a cli.Args) ([]string, error) {
	var err error

	if err = defineProjectMainContainer(cfg, d, c.GetContainerList()); err != nil {
		return []string{}, err
	}

	var initArgs = []string{"exec"}

	composerArgs := []string{"-it", cfg.GetProjectMainContainer(), "composer"}

	if c.GetCallType() == "memory" || a.Get(0) == "m" {
		var phpLocation, composerLocation string

		getCommandLocation := c.GetCommandLocaton()

		if phpLocation, err = getCommandLocation(cfg.GetProjectMainContainer(), "php"); err != nil {
			return nil, err
		}

		if composerLocation, err = getCommandLocation(cfg.GetProjectMainContainer(), "composer"); err != nil {
			return nil, err
		}

		composerArgs = []string{"-i", cfg.GetProjectMainContainer(), phpLocation, "-d", "memory_limit=-1", composerLocation}
	}

	if len(c.GetComposerCommand()) > 0 {
		composerArgs = append(composerArgs, c.GetComposerCommand())
	}

	initArgs = append(initArgs, composerArgs...)

	extraInitArgs := a.Slice()

	if a.Get(0) == "m" {
		extraInitArgs = a.Tail()
	}

	args := append(initArgs, extraInitArgs...)

	return args, nil
}
