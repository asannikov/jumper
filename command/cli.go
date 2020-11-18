package command

import (
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2" // imports as package "cli"
)

// CallCliCommand generates CLI command (docker-compose exec phpfpm "$@")
func CallCliCommand(containerPhp string) *cli.Command {

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
				// https://medium.com/@ssttehrani/containers-from-scratch-with-golang-5276576f9909
				// https://phase2.github.io/devtools/common-tasks/ssh-into-a-container/
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
