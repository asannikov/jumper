package command

import (
	"errors"

	"github.com/urfave/cli/v2"
)

type shellConfig interface {
	GetShell() string
	SaveShellCommand(string) error
}

type shellDialog interface {
	DockerShell() (int, string, error)
}

type shellOptions interface {
	GetInitFuntion() func(bool) string
}

// ShellCommand changes shell type
func ShellCommand(cfg shellConfig, d shellDialog, options shellOptions) *cli.Command {
	initf := options.GetInitFuntion()
	return &cli.Command{
		Name:            "shell",
		Usage:           "Change shell type for a project",
		SkipFlagParsing: false,
		Action: func(c *cli.Context) (err error) {
			initf(true)

			if err = defineShellType(cfg, d); err != nil {
				return err
			}

			return nil
		},
	}
}

func defineShellType(cfg shellConfig, d shellDialog) (err error) {
	var path string

	if _, path, err = d.DockerShell(); err != nil {
		return err
	}

	if path == "" {
		return errors.New("Something goes wrong. Shell was not set")
	}

	return cfg.SaveShellCommand(path)
}
