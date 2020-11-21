package command

import (
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2" // imports as package "cli"
)

// CallCliCommand generates CLI command (docker-compose exec phpfpm "$@")
func CallCliCommand(cp func() string) *cli.Command {

	containerPhp := cp()

	cmd := cli.Command{
		Name:    "cli",
		Aliases: []string{"c"},
		Usage:   "Run cli",
		Action: func(c *cli.Context) error {
			dockerComposeCliCommand := c.Args().Get(0)

			if dockerComposeCliCommand == "" {
				log.Println("Please specify a CLI command (ex. ls)")
			} else {
				var binary = "docker"
				var initArgs = []string{"exec", "-it", containerPhp, dockerComposeCliCommand}

				extraInitArgs := []string{}

				args := append(initArgs, extraInitArgs...)

				log.Println(args)
				cmd := exec.Command(binary, args...)

				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				cmd.Run()

			}

			return nil
		},
	}

	return &cmd
}
