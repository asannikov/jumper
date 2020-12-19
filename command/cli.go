package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2" // imports as package "cli"
)

// CallCliCommand generates CLI command (docker-compose exec phpfpm "$@")
func CallCliCommand(initf func(), cfg projectConfig, d dialog, containerlist []string) *cli.Command {
	cmd := cli.Command{
		Name:    "cli",
		Aliases: []string{"c"},
		Usage:   "Runs cli in container: {docker exec -it main_container} [bash command] [custom parameters]",
		Action: func(c *cli.Context) (err error) {
			initf()

			var args []string

			dockerComposeCliCommand := c.Args().Get(0)

			if dockerComposeCliCommand == "" {
				return errors.New("Please specify a CLI command (ie. ls)")
			}

			if args, err = cliCommandHandle(cfg, d, containerlist, dockerComposeCliCommand, c.Args()); err != nil {
				return err
			}

			fmt.Printf("\n command: %s\n\n", "docker "+strings.Join(args, " "))

			cmd := exec.Command("docker", args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		},
	}

	return &cmd
}

// CallBashCommand generates Bash command (docker-compose exec phpfpm bash)
func CallBashCommand(initf func(), cfg projectConfig, d dialog, containerlist []string) *cli.Command {
	cmd := cli.Command{
		Name:    "bash",
		Aliases: []string{"b"},
		Usage:   "Runs cli in container: {docker exec -it main_container bash} [custom parameters]",
		Action: func(c *cli.Context) (err error) {
			initf()

			var args []string

			if args, err = cliCommandHandle(cfg, d, containerlist, "bash", c.Args()); err != nil {
				return err
			}

			fmt.Printf("\n command: %s\n\n", "docker "+strings.Join(args, " "))

			cmd := exec.Command("docker", args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		},
	}

	return &cmd
}

func cliCommandHandle(cfg projectConfig, d dialog, containerlist []string, command string, a cli.Args) ([]string, error) {
	var err error

	if err = defineProjectMainContainer(cfg, d, containerlist); err != nil {
		return []string{}, err
	}

	var initArgs = []string{"exec", "-it", cfg.GetProjectMainContainer(), command}
	extraInitArgs := a.Tail()
	return append(initArgs, extraInitArgs...), nil
}
