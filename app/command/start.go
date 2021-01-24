package command

import (
	"errors"
	"strings"

	"github.com/urfave/cli/v2"
)

type defineStartCommandProjectConfig interface {
	GetStartCommand() string
	SaveStartCommandToProjectConfig(string) error
}

type defineStartCommandDialog interface {
	StartCommand() (string, error)
}

func defineStartCommand(cfg defineStartCommandProjectConfig, d defineStartCommandDialog, containerlist []string) (err error) {
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

type runStartProjectProjectConfig interface {
	GetStartCommand() string
}

type runStartProjectOptions interface {
	GetExecCommand() func(ExecOptions, *cli.App) error
}

func runStartProject(c *cli.Context, cfg runStartProjectProjectConfig, args []string, options runStartProjectOptions) error {
	commandSlice := strings.Split(cfg.GetStartCommand(), " ")
	execCommand := options.GetExecCommand()

	var binary = commandSlice[0]
	var initArgs = commandSlice[1:]

	extraInitArgs := c.Args().Slice()

	args = append(initArgs, args...)
	args = append(args, extraInitArgs...)

	eo := ExecOptions{
		command: binary,
		args:    args,
		tty:     true,
		detach:  true,
	}

	return execCommand(eo, c.App)
}

type callStartProjectBasicProjectConfig interface {
	GetProjectMainContainer() string
	SaveContainerNameToProjectConfig(string) error
	GetStartCommand() string
	SaveStartCommandToProjectConfig(string) error
}

type callStartProjectBasicDialog interface {
	SetMainContaner([]string) (int, string, error)
	StartCommand() (string, error)
}

type startProjectOptions interface {
	GetInitFuntion() func(bool) string
	GetContainerList() ([]string, error)
	GetExecCommand() func(ExecOptions, *cli.App) error
}

// CallStartProjectBasic runs docker project
func CallStartProjectBasic(cfg callStartProjectBasicProjectConfig, d callStartProjectBasicDialog, options startProjectOptions) *cli.Command {
	initf := options.GetInitFuntion()

	cmd := cli.Command{
		Name:            "start",
		Aliases:         []string{"st"},
		Usage:           `Runs defined command: {docker-compose -f docker-compose.yml up} [custom parameters]`,
		Description:     `It's possible to use any custom parameters coming after "up"`,
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf(true)

			var cl []string

			if cl, err = options.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			if err = defineStartCommand(cfg, d, cl); err != nil {
				return err
			}

			return runStartProject(c, cfg, []string{}, options)
		},
	}

	return &cmd
}

// CallStartProjectForceRecreate runs docker project
func CallStartProjectForceRecreate(cfg callStartProjectBasicProjectConfig, d callStartProjectBasicDialog, options startProjectOptions) *cli.Command {
	initf := options.GetInitFuntion()

	cmd := cli.Command{
		Name:    "start:force",
		Aliases: []string{"s:f"},
		Usage:   `Runs defined command: {docker-compose -f docker-compose.yml up --force-recreat} [custom parameters]`,
		Description: `
		--force-recreate - Recreate containers even if their configuration and image haven't changed
		It's possible to use any custom parameters coming after "up"`,
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf(true)

			var cl []string

			if cl, err = options.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			if err = defineStartCommand(cfg, d, cl); err != nil {
				return err
			}

			return runStartProject(c, cfg, []string{"--force-recreate"}, options)
		},
	}

	return &cmd
}

// CallStartProjectOrphans runs docker project
func CallStartProjectOrphans(cfg callStartProjectBasicProjectConfig, d callStartProjectBasicDialog, options startProjectOptions) *cli.Command {
	initf := options.GetInitFuntion()

	cmd := cli.Command{
		Name:    "start:orphans",
		Aliases: []string{"s:o"},
		Usage:   `Runs defined command: {docker-compose -f docker-compose.yml up --remove-orphans} [custom parameters]`,
		Description: `
		--remove-orphans - Remove containers for services not defined in the Compose file
		It's possible to use any custom parameters coming after "up"`,
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf(true)

			var cl []string

			if cl, err = options.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			if err = defineStartCommand(cfg, d, cl); err != nil {
				return err
			}

			return runStartProject(c, cfg, []string{"--remove-orphans"}, options)
		},
	}

	return &cmd
}

// CallStartProjectForceOrphans runs docker project
func CallStartProjectForceOrphans(cfg callStartProjectBasicProjectConfig, d callStartProjectBasicDialog, options startProjectOptions) *cli.Command {
	initf := options.GetInitFuntion()

	cmd := cli.Command{
		Name:    "start:force-orphans",
		Aliases: []string{"s:fo"},
		Usage:   `Runs defined command: {docker-compose -f docker-compose.yml up --force-recreate --remove-orphans} [custom parameters]`,
		Description: `
		--force-recreate - Recreate containers even if their configuration and image haven't changed
		--remove-orphans - Remove containers for services not defined in the Compose file
		It's possible to use any custom parameters coming after "up"`,
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf(true)

			var cl []string

			if cl, err = options.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			if err = defineStartCommand(cfg, d, cl); err != nil {
				return err
			}

			return runStartProject(c, cfg, []string{"--force-recreate", "--remove-orphans"}, options)
		},
	}

	return &cmd
}

type callStartMainContainerProjectConfig interface {
	GetProjectMainContainer() string
	SaveContainerNameToProjectConfig(string) error
}

// CallStartMainContainer runs docker main container
func CallStartMainContainer(cfg callStartMainContainerProjectConfig, d callStartProjectBasicDialog, options startProjectOptions) *cli.Command {
	initf := options.GetInitFuntion()
	execCommand := options.GetExecCommand()

	cmd := cli.Command{
		Name:    "start:maincontainer",
		Aliases: []string{"startmc"},
		Usage:   `Runs defined command: {docker start main_container}`,
		Action: func(c *cli.Context) (err error) {
			initf(true)

			var cl []string

			if cl, err = options.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			args := []string{"start", cfg.GetProjectMainContainer()}

			eo := ExecOptions{
				command: "docker",
				args:    args,
				tty:     true,
				detach:  true,
			}

			return execCommand(eo, c.App)
		},
	}

	return &cmd
}

type restartMainContainerProjectConfig interface {
	GetProjectMainContainer() string
}

type restartMainContainerOptions interface {
	GetExecCommand() func(ExecOptions, *cli.App) error
}

func restartMainContainer(cfg restartMainContainerProjectConfig, options restartMainContainerOptions, a *cli.App) error {
	execCommand := options.GetExecCommand()

	args := []string{"stop", cfg.GetProjectMainContainer()}

	eo := ExecOptions{
		command: "docker",
		args:    args,
		tty:     true,
		detach:  true,
	}

	if err := execCommand(eo, a); err != nil {
		return err
	}

	eo = ExecOptions{
		command: "docker",
		args:    []string{"start", cfg.GetProjectMainContainer()},
		tty:     true,
		detach:  true,
	}

	return execCommand(eo, a)
}

type restartProjectOptions interface {
	GetInitFuntion() func(bool) string
	GetContainerList() ([]string, error)
	GetExecCommand() func(ExecOptions, *cli.App) error
	GetDockerStatus() bool
}

// CallRestartMainContainer restarts docker main container
func CallRestartMainContainer(cfg callStartMainContainerProjectConfig, d callStartProjectBasicDialog, options restartProjectOptions) *cli.Command {
	initf := options.GetInitFuntion()
	dockerStatus := options.GetDockerStatus()

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

			if cl, err = options.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			return restartMainContainer(cfg, options, c.App)
		},
	}

	return &cmd
}

type callStartContainersOptions interface {
	GetExecCommand() func(ExecOptions, *cli.App) error
	GetInitFuntion() func(bool) string
}

// CallStartContainers runs docker custom container
func CallStartContainers(options callStartContainersOptions) *cli.Command {
	initf := options.GetInitFuntion()
	execCommand := options.GetExecCommand()

	cmd := cli.Command{
		Name:    "start:containers",
		Aliases: []string{"startc"},
		Usage:   `Runs defined command: {docker start} [container]`,
		Action: func(c *cli.Context) (err error) {
			initf(true)

			eo := ExecOptions{
				command: "docker",
				args:    append([]string{"start"}, c.Args().Slice()...),
				tty:     true,
				detach:  true,
			}

			return execCommand(eo, c.App)
		},
	}

	return &cmd
}

type callRestartContainersOptions interface {
	GetExecCommand() func(ExecOptions, *cli.App) error
	GetInitFuntion() func(bool) string
	GetDockerStatus() bool
}

// CallRestartContainers restart docker custom containers
func CallRestartContainers(options callRestartContainersOptions) *cli.Command {
	initf := options.GetInitFuntion()
	dockerStatus := options.GetDockerStatus()
	execCommand := options.GetExecCommand()

	cmd := cli.Command{
		Name:    "restart:containers",
		Aliases: []string{"rc"},
		Usage:   `Runs defined command: {docker start} [container]`,
		Action: func(c *cli.Context) (err error) {
			if !dockerStatus {
				return errors.New("Docker is not running")
			}

			initf(true)

			eo := ExecOptions{
				command: "docker",
				args:    append([]string{"stop"}, c.Args().Slice()...),
				tty:     true,
				detach:  true,
			}

			if err := execCommand(eo, c.App); err != nil {
				return err
			}

			eo = ExecOptions{
				command: "docker",
				args:    append([]string{"start"}, c.Args().Slice()...),
				tty:     true,
				detach:  true,
			}
			return execCommand(eo, c.App)
		},
	}

	return &cmd
}
