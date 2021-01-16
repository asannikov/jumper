package command

import (
	"errors"

	"github.com/urfave/cli/v2" // imports as package "cli"
)

type commandHandler func(chc commandHandleProjectConfig) string

type cliCommand struct {
	usage   map[string]string
	aliases map[string]string
	args    map[string][]string
	command map[string]func(commandHandleProjectConfig) string
}

func (cli *cliCommand) GetCommand(cmd string, cfg commandHandleProjectConfig) string {
	command := cli.command[cmd]
	return command(cfg)
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

type commandHandleProjectConfig interface {
	GetShell() string
}

type callCliCommandOptions interface {
	GetInitFuntion() func(bool) string
	GetContainerList() ([]string, error)
	GetExecCommand() func(string, []string, *cli.App) error
}

// CallCliCommand calls a range of differnt cli commands
func CallCliCommand(commandName string, cfg cliCommandHandleProjectConfig, d callCliCommandDialog, options callCliCommandOptions) *cli.Command {
	initf := options.GetInitFuntion()
	execCommand := options.GetExecCommand()

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
		command: map[string]func(commandHandleProjectConfig) string{
			"cli": func(c commandHandleProjectConfig) string {
				return ""
			},
			"sh": func(c commandHandleProjectConfig) string {
				return c.GetShell()
			},
			"clinotty": func(c commandHandleProjectConfig) string {
				return ""
			},
			"cliroot": func(c commandHandleProjectConfig) string {
				return ""
			},
			"clirootnotty": func(c commandHandleProjectConfig) string {
				return ""
			},
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

			if args, err = cliCommandHandle(commandName, cfg, d, clic, options, c.Args()); err != nil {
				return err
			}

			return execCommand("docker", args, c.App)
		},
	}
}

type cliCommandInterface interface {
	GetCommand(string, commandHandleProjectConfig) string
	GetArgs() map[string][]string
}

type cliCommandHandleDialog interface {
	SetMainContaner([]string) (int, string, error)
}

func cliCommandHandle(index string, cfg cliCommandHandleProjectConfig, d cliCommandHandleDialog, c cliCommandInterface, options containerlist, a cli.Args) ([]string, error) {
	var err error
	var cl []string

	if cl, err = options.GetContainerList(); err != nil {
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

	if c.GetCommand(index, cfg) != "" {
		initArgs = append(initArgs, c.GetCommand(index, cfg))
	}

	if c.GetCommand(index, cfg) == "" && a.Get(0) == "" {
		return []string{}, errors.New("Please specify a CLI command (ex. ls)")
	}

	return append(initArgs, a.Slice()...), nil
}
