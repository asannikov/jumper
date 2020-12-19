package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2" // imports as package "cli"
)

type cliCommand struct {
	usage         map[string]string
	aliases       map[string]string
	args          map[string][]string
	containerList []string
	command       map[string]string
}

func (cli *cliCommand) GetContainerList() []string {
	return cli.containerList
}

func (cli *cliCommand) GetCommand(cmd string) string {
	return cli.command[cmd]
}

func (cli *cliCommand) GetArgs() map[string][]string {
	return cli.args
}

// CallCliCommand calls a range of differnt cli commands
func CallCliCommand(commandName string, initf func(), cfg projectConfig, d dialog, containerlist []string) *cli.Command {
	clic := &cliCommand{
		usage: map[string]string{
			"cli":          "Runs cli command in conatiner: {docker exec main_conatain} [command] [custom parameters]",
			"bash":         "Runs cli bash command in conatiner: {docker exec main_conatain bash} [custom parameters]",
			"clinotty":     "Runs command {docker exec -t main_container} [command] [custom parameters]",
			"cliroot":      "Runs command {docker exec -u root main_container} [command] [custom parameters]",
			"clirootnotty": "Runs command {docker exec -u root -t main_container} [command] [custom parameters]",
		},
		aliases: map[string]string{
			"cli":          "c",
			"bash":         "b",
			"clinotty":     "cnt",
			"cliroot":      "cr",
			"clirootnotty": "crnt",
		},
		args: map[string][]string{
			"cli":          []string{"-it"},
			"bash":         []string{"-it"},
			"cliroot":      []string{"-u", "root", "-it"},
			"clinotty":     []string{"-i"},
			"clirootnotty": []string{"-u", "root", "-i"},
		},
		containerList: containerlist,
		command: map[string]string{
			"cli":          "",
			"bash":         "bash",
			"clinotty":     "",
			"cliroot":      "",
			"clirootnotty": "",
		},
	}

	return &cli.Command{
		Name:            commandName,
		Aliases:         []string{clic.aliases[commandName]},
		Usage:           clic.usage[commandName],
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf()

			var args []string

			if args, err = cliCommandHandle(commandName, cfg, d, clic, c.Args()); err != nil {
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

type cliCommandInterface interface {
	GetContainerList() []string
	GetCommand(string) string
	GetArgs() map[string][]string
}

func cliCommandHandle(index string, cfg projectConfig, d dialog, c cliCommandInterface, a cli.Args) ([]string, error) {
	var err error

	if err = defineProjectMainContainer(cfg, d, c.GetContainerList()); err != nil {
		return []string{}, err
	}

	var initArgs = []string{"exec"}

	extraInitArgs := c.GetArgs()[index]

	if len(extraInitArgs) > 0 {
		initArgs = append(initArgs, extraInitArgs...)
	}

	initArgs = append(initArgs, cfg.GetProjectMainContainer())

	if c.GetCommand(index) != "" {
		initArgs = append(initArgs, c.GetCommand(index))
	}

	if c.GetCommand(index) == "" && a.Get(0) == "" {
		return []string{}, errors.New("Please specify a CLI command (ex. ls)")
	}

	return append(initArgs, a.Slice()...), nil
}
