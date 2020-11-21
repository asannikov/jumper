package command

import (
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2" // imports as package "cli"
)

// CallBashCommand generates Bash command (docker-compose exec phpfpm bash)
func CallBashCommand(cp func() string) *cli.Command {

	containerPhp := cp()

	cmd := cli.Command{
		Name:    "bash",
		Aliases: []string{"b"},
		Usage:   "Run bash",
		Action: func(c *cli.Context) error {

			var binary = "docker"
			var initArgs = []string{"exec", "-it", containerPhp, "bash"}

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
