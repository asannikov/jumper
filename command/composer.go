package command

import (
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2" // imports as package "cli"
)

// CallComposerCommand generates Composer command (docker-compose exec phpfpm composer)
// @todo do composer install/update with no memory limit
func CallComposerCommand(containerPhp string) *cli.Command {

	cmd := cli.Command{
		Name:    "composer",
		Aliases: []string{"cmp"},
		Usage:   "Run composer",
		Action: func(c *cli.Context) error {
			var binary = "docker"
			var initArgs = []string{"exec", "-it", containerPhp, "composer"}

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

			return nil
		},
	}

	return &cmd
}
