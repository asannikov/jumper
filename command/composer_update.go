package command

import (
	"github.com/urfave/cli/v2" // imports as package "cli"
	"os"
	"os/exec"
)

// CallComposerUpdateCommand generates Composer command (docker-compose exec phpfpm composer)
// @todo do composer install/update with no memory limit
func CallComposerUpdateCommand(cp func() (string, error)) *cli.Command {

	cmd := cli.Command{
		Name:            "composer:update",
		Aliases:         []string{"cmpu"},
		Usage:           "Run composer update",
		SkipFlagParsing: true,
		Action: func(c *cli.Context) error {
			containerPhp, _ := cp()

			if c.Args().Get(0) == "m" {
				return explictComposerUpdate(containerPhp, c.Args().Tail())
			}

			var binary = "docker"
			var initArgs = []string{"exec", "-it", containerPhp, "composer", "update"}

			extraInitArgs := c.Args().Slice()

			args := append(initArgs, extraInitArgs...)
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

// CallComposerUpdateMemoryCommand generates Composer command (docker-compose exec phpfpm php -d memory_limit composer update)
func CallComposerUpdateMemoryCommand(cp func() (string, error)) *cli.Command {

	cmd := cli.Command{
		Name:            "composer:update:memory",
		Aliases:         []string{"cmpum"},
		Usage:           "Run php -d memory_limit=-1 composer update",
		SkipFlagParsing: true,
		Action: func(c *cli.Context) error {
			containerPhp, _ := cp()
			return explictComposerUpdate(containerPhp, c.Args().Slice())
		},
	}

	return &cmd
}

func explictComposerUpdate(phpContainer string, extraInitArgs []string) error {
	cmp := Composer{}
	phpContainerName, composerContainerName, err := cmp.GetContainerNames(phpContainer)

	if err == nil {
		var binary = "docker"
		var initArgs = []string{"exec", "-i", phpContainer, phpContainerName, "-d", "memory_limit=-1", composerContainerName, "update"}

		args := append(initArgs, extraInitArgs...)

		cmd := exec.Command(binary, args...)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.Run()
	}

	return err
}
