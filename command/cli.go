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
	usage   map[string]string
	aliases map[string]string
	args    map[string][]string
	command map[string]string
}

func (cli *cliCommand) GetCommand(cmd string) string {
	return cli.command[cmd]
}

func (cli *cliCommand) GetArgs() map[string][]string {
	return cli.args
}

type cliCommandHandleProjectConfig interface {
	GetProjectMainContainer() string
	SaveContainerNameToProjectConfig(string) error
	GetShell() string
}

type callCliCommandDialog interface {
	SetMainContaner([]string) (int, string, error)
}

// CallCliCommand calls a range of differnt cli commands
func CallCliCommand(commandName string, initf func(bool) string, cfg cliCommandHandleProjectConfig, d callCliCommandDialog, cl containerlist) *cli.Command {
	clic := &cliCommand{
		usage: map[string]string{
			"cli":          "Runs cli command in conatiner: {docker exec main_conatain} [command] [custom parameters]",
			"sh":           "Runs cli sh command in conatiner: {docker exec main_conatain {shell_type}} [custom parameters]",
			"clinotty":     "Runs command {docker exec -t main_container} [command] [custom parameters]",
			"cliroot":      "Runs command {docker exec -u root main_container} [command] [custom parameters]",
			"clirootnotty": "Runs command {docker exec -u root -t main_container} [command] [custom parameters]",
		},
		aliases: map[string]string{
			"cli":          "c",
			"sh":           "sh",
			"clinotty":     "cnt",
			"cliroot":      "cr",
			"clirootnotty": "crnt",
		},
		args: map[string][]string{
			"cli":          []string{"-it"},
			"sh":           []string{"-it"},
			"cliroot":      []string{"-u", "root", "-it"},
			"clinotty":     []string{"-i"},
			"clirootnotty": []string{"-u", "root", "-i"},
		},
		command: map[string]string{
			"cli":          "",
			"sh":           cfg.GetShell(),
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
			initf(true)

			var args []string

			if args, err = cliCommandHandle(commandName, cfg, d, clic, cl, c.Args()); err != nil {
				return err
			}

			fmt.Printf("\ncommand: %s\n\n", "docker "+strings.Join(args, " "))

			cmd := exec.Command("docker", args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		},
	}
}

type cliCommandInterface interface {
	GetCommand(string) string
	GetArgs() map[string][]string
}

type cliCommandHandleDialog interface {
	SetMainContaner([]string) (int, string, error)
}

func cliCommandHandle(index string, cfg cliCommandHandleProjectConfig, d cliCommandHandleDialog, c cliCommandInterface, clist containerlist, a cli.Args) ([]string, error) {
	var err error
	var cl []string

	if cl, err = clist.GetContainerList(); err != nil {
		return []string{}, err
	}

	if err = defineProjectMainContainer(cfg, d, cl); err != nil {
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
