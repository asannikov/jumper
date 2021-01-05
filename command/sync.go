package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

type sync struct {
	usage       map[string]string
	aliases     map[string]string
	description map[string]string
}

func getSyncPath(path string) string {
	if path == "--all" || path == "."+string(os.PathSeparator) || path == string(os.PathSeparator)+"." || path == "." {
		return "/./"
	}

	return string(os.PathSeparator) + strings.Trim(path, string(os.PathSeparator))
}

func getSyncArgs(cfg projectConfig, direction string, syncPath string, projectRoot string) []string {
	projectRoot = strings.TrimRight(projectRoot, string(os.PathSeparator))

	args := []string{"cp", projectRoot + syncPath, cfg.GetProjectMainContainer() + ":" + strings.TrimRight(cfg.GetProjectDockerPath(), string(os.PathSeparator)) + syncPath}

	if direction == "copyto" {
		args = []string{"cp", cfg.GetProjectMainContainer() + ":" + strings.TrimRight(cfg.GetProjectDockerPath(), string(os.PathSeparator)) + syncPath, projectRoot + syncPath}
	}

	return args
}

//SyncCommand does the syncronization between container and project
func SyncCommand(direction string, initf func(bool) string, dockerStatus bool, cfg projectConfig, d dialog, clist containerlist) *cli.Command {

	s := &sync{
		usage: map[string]string{
			"copyto":   "Sync local -> docker container, set related path, ie `vendor/folder/` for syncing as a parameter, or use --all to sync all project",
			"copyfrom": "Sync docker container -> local, set related path, ie `vendor/folder/` for syncing as a parameter, or use --all to sync all project",
		},
		aliases: map[string]string{
			"copyto":   "cpt",
			"copyfrom": "cpf",
		},
		description: map[string]string{
			"copyto":   "phpContainer is taken from project config file",
			"copyfrom": "phpContainer is taken from project config file, php and composer commands will be found automatically",
		},
	}

	return &cli.Command{
		Name:            direction,
		Aliases:         []string{s.aliases[direction]},
		Usage:           s.usage[direction],
		Description:     s.description[direction],
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			syncPath := c.Args().First()

			if syncPath == "" {
				return errors.New("Please, specify the path you want to sync")
			}

			currentPath := initf(true)
			syncPath = getSyncPath(syncPath)

			var cl []string

			if cl, err = clist.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			if err = defineProjectDockerPath(cfg, d, "/var/www/html/"); err != nil {
				return err
			}

			args := getSyncArgs(cfg, direction, syncPath, currentPath)

			cmd := exec.Command("docker", args...)

			fmt.Printf("\ncommand: %s\n\n", "docker "+strings.Join(args, " "))

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err = cmd.Run(); err != nil {
				return err
			}

			if direction == "copyto" {
				fmt.Printf("Completed copying %s files from host to container %s \n", syncPath, cfg.GetProjectMainContainer())
			} else {
				fmt.Printf("Completed copying %s from container %s to host\n", syncPath, cfg.GetProjectMainContainer())
			}

			return nil
		},
	}
}
