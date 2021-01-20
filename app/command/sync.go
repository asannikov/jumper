package command

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

const syncCopyFrom = "copyfrom"
const syncCopyTo = "copyto"

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

type syncProjectConfig interface {
	GetProjectMainContainer() string
	GetProjectDockerPath() string
	SaveContainerNameToProjectConfig(string) error
	SaveDockerProjectPath(string) error
}

func getSyncArgs(cfg syncProjectConfig, direction string, syncPath string, projectRoot string) []string {
	projectRoot = strings.TrimRight(projectRoot, string(os.PathSeparator))

	targetPath := filepath.Dir(syncPath)

	args := []string{"cp", projectRoot + syncPath, cfg.GetProjectMainContainer() + ":" + strings.TrimRight(cfg.GetProjectDockerPath(), string(os.PathSeparator)) + targetPath}

	if direction == syncCopyFrom {
		args = []string{"cp", cfg.GetProjectMainContainer() + ":" + strings.TrimRight(cfg.GetProjectDockerPath(), string(os.PathSeparator)) + syncPath, projectRoot + targetPath}
	}

	return args
}

type syncCommandDialog interface {
	SetMainContaner([]string) (int, string, error)
	DockerProjectPath(string) (string, error)
}

type syncOptions interface {
	GetExecCommand() func(string, []string, *cli.App) error
	GetInitFuntion() func(bool) string
	GetContainerList() ([]string, error)
}

//SyncCommand does the syncronization between container and project
func SyncCommand(direction string, cfg syncProjectConfig, d syncCommandDialog, options syncOptions) *cli.Command {
	execCommand := options.GetExecCommand()
	initf := options.GetInitFuntion()

	s := &sync{
		usage: map[string]string{
			syncCopyTo:   "Sync local -> docker container, set related path, ie `vendor/folder/` for syncing as a parameter, or use --all to sync all project",
			syncCopyFrom: "Sync docker container -> local, set related path, ie `vendor/folder/` for syncing as a parameter, or use --all to sync all project",
		},
		aliases: map[string]string{
			syncCopyTo:   "cpt",
			syncCopyFrom: "cpf",
		},
		description: map[string]string{
			syncCopyTo:   "Works only for defined main container. Keep in mind that `docker cp` create only the top folder of the path if all nodes of the path do not exist. For such case use -f flag. It creates all folders recursively.",
			syncCopyFrom: "phpContainer is taken from project config file",
		},
	}

	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "Force create directory for file if it does not exist",
		},
	}

	return &cli.Command{
		Name:            direction,
		Aliases:         []string{s.aliases[direction]},
		Usage:           s.usage[direction],
		Description:     s.description[direction],
		SkipFlagParsing: false,
		Flags:           flags,
		Action: func(c *cli.Context) (err error) {
			syncPath := c.Args().First()

			if syncPath == "" {
				return errors.New("Please, specify the path you want to sync")
			}

			currentPath := initf(true)
			syncPath = getSyncPath(syncPath)

			var cl []string

			if cl, err = options.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			if err = defineProjectDockerPath(cfg, d, "/var/www/html/"); err != nil {
				return err
			}

			args := getSyncArgs(cfg, direction, syncPath, currentPath)

			if direction == syncCopyFrom && c.Bool("f") == true {
				fmt.Printf("Path %s was created", args[2]+filepath.Base(syncPath))
				err = os.MkdirAll(args[2]+filepath.Base(syncPath), os.ModePerm)
			}

			if direction == syncCopyTo && c.Bool("f") == true {
				fmt.Printf("Path %s was created", cfg.GetProjectDockerPath()+syncPath)
				// @todo - use go client logic, branch 32-add-force-copy-flag
				err = execCommand("docker", []string{"exec", cfg.GetProjectMainContainer(), "mkdir", "-p", cfg.GetProjectDockerPath() + syncPath}, c.App)
			}

			if err != nil {
				return err
			}

			if err = execCommand("docker", args, c.App); err != nil {
				return err
			}

			if direction == syncCopyTo {
				fmt.Printf("Completed copying %s files from host to container %s \n", syncPath, cfg.GetProjectMainContainer())
			} else {
				fmt.Printf("Completed copying %s from container %s to host\n", syncPath, cfg.GetProjectMainContainer())
			}

			return nil
		},
	}
}
