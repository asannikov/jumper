package command

import (
	"os"
	"os/exec"

	"github.com/urfave/cli/v2" // imports as package "cli"
)

// CallComposerCommand generates Composer command (docker-compose exec phpfpm composer)
func CallComposerCommand(cp func() (string, error)) *cli.Command {

	cmd := cli.Command{
		Name:            "composer",
		Aliases:         []string{"cmp"},
		Usage:           "Run composer",
		SkipFlagParsing: true,
		Action: func(c *cli.Context) error {
			containerPhp, _ := cp()

			var binary = "docker"
			var initArgs = []string{"exec", "-it", containerPhp, "composer"}

			extraInitArgs := c.Args().Slice()

			args := append(initArgs, extraInitArgs...)

			cmd := exec.Command(binary, args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			cmd.Run()

			return nil
		},
	}

	return &cmd
}
