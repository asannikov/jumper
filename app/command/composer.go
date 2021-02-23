package command

import (
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

type callComposerCommandProjectConfig interface {
	GetProjectDockerPath() string
	GetProjectMainContainer() string
	SaveContainerNameToProjectConfig(string) error
	GetCommandInactveStatus(string) bool
	SaveContainerUserToProjectConfig(string) error
	GetMainContainerUser() string
}

type callComposerCommandDialog interface {
	SetMainContaner([]string) (int, string, error)
	SetMainContanerUser() (string, error)
}

type callComposerCommandOptions interface {
	GetInitFunction() func(bool) string
	GetContainerList() ([]string, error)
	GetExecCommand() func(ExecOptions, *cli.App) error
	GetCommandLocation() func(string, string) (string, error)
}

// CallComposerCommand generates composer commands
// https://medium.com/@ssttehrani/containers-from-scratch-with-golang-5276576f9909
// https://phase2.github.io/devtools/common-tasks/ssh-into-a-container/
func CallComposerCommand(composercommand string, cfg callComposerCommandProjectConfig, d callComposerCommandDialog, options callComposerCommandOptions) *cli.Command {
	index, calltype, dockercmd := parseCommand(composercommand)

	initf := options.GetInitFunction()
	execCommand := options.GetExecCommand()

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
		locaton: options.GetCommandLocation(),
		ctype:   calltype,
		command: dockercmd,
	}

	return &cli.Command{
		Name:            composercommand,
		Aliases:         []string{cmp.aliases[index]},
		Usage:           cmp.usage[index],
		Description:     cmp.description[index],
		Hidden:          cfg.GetCommandInactveStatus("composer"),
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf(true)

			var args []string

			if args, err = composerHandle(cfg, d, cmp, options, c.Args()); err != nil {
				return err
			}

			eo := ExecOptions{
				command: "docker",
				args:    args,
				tty:     true,
				detach:  true,
			}

			return execCommand(eo, c.App)
		},
	}
}

type composerInterface interface {
	GetCommandLocaton() func(string, string) (string, error)
	GetCallType() string
	GetComposerCommand() string
}

type composerHandleProjectConfig interface {
	GetProjectMainContainer() string
	SaveContainerNameToProjectConfig(string) error
	GetMainContainerUser() string
	SaveContainerUserToProjectConfig(string) error
}

type composerHandleDialog interface {
	SetMainContaner([]string) (int, string, error)
	SetMainContanerUser() (string, error)
}

func composerHandle(cfg composerHandleProjectConfig, d composerHandleDialog, c composerInterface, options containerlist, a cli.Args) ([]string, error) {
	var err error
	var cl []string

	if cl, err = options.GetContainerList(); err != nil {
		return []string{}, err
	}

	if err = defineProjectMainContainer(cfg, d, cl); err != nil {
		return []string{}, err
	}

	if err = defineProjectMainContainerUser(cfg, d); err != nil {
		return []string{}, err
	}

	var initArgs = []string{"exec"}

	composerArgs := []string{"-it", cfg.GetProjectMainContainer(), "composer"}

	if cfg.GetMainContainerUser() != "root" && cfg.GetMainContainerUser() != "" {
		composerArgs = []string{"-it", "-u", cfg.GetMainContainerUser(), cfg.GetProjectMainContainer(), "composer"}
	}

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

		if cfg.GetMainContainerUser() != "root" && cfg.GetMainContainerUser() != "" {
			composerArgs = []string{"-i", "-u", cfg.GetMainContainerUser(), cfg.GetProjectMainContainer(), phpLocation, "-d", "memory_limit=-1", composerLocation}
		}
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
