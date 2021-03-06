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
	GetMainContainerUser() string
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
	GetExecCommand() func(ExecOptions, *cli.App) error
	GetInitFunction() func(bool) string
	GetContainerList() ([]string, error)
	GetCopyTo(string, string, string) error
	RunNativeExec(ExecOptions, *cli.App) error
	DirExists(string) (bool, error)
	MkdirAll(string, os.FileMode) error
}

// SyncCommand does the syncronization between container and project
func SyncCommand(direction string, cfg syncProjectConfig, d syncCommandDialog, options syncOptions) *cli.Command {
	execCommand := options.GetExecCommand()
	initf := options.GetInitFunction()

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
			syncCopyFrom: "Works only for defined main container. Keep in mind that `docker cp` create only the top folder of the path if all nodes of the path do not exist. For such case use -f flag. It creates all folders recursively.",
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
				if err = options.MkdirAll(args[2]+string(os.PathSeparator)+filepath.Base(syncPath), os.ModePerm); err == nil {
					fmt.Printf("Path %s was created", args[2]+string(os.PathSeparator)+filepath.Base(syncPath))
				}
			}

			if direction == syncCopyTo && c.Bool("f") == true {
				eo := ExecOptions{
					command: "mkdir",
					args:    []string{"-p", cfg.GetProjectDockerPath() + strings.TrimLeft(syncPath, string(os.PathSeparator))},
					tty:     true,
					detach:  true,
					user:    "root",
				}

				status := false

				status, err = options.DirExists(args[1])
				if err != nil || status == false {
					fmt.Printf("Source directory %s does not exist\n", args[1])
					return err
				}

				if err = options.RunNativeExec(eo, c.App); err == nil {
					fmt.Printf("Path %s was created", args[2]+string(os.PathSeparator)+filepath.Base(syncPath))
				}
			}

			if err != nil {
				return err
			}

			eo := ExecOptions{
				command: "docker",
				args:    args,
				tty:     true,
				detach:  true,
			}

			if err = execCommand(eo, c.App); err != nil {
				return err
			}

			if direction == syncCopyTo {
				fmt.Printf("Completed copying %s files from host to container %s \n", syncPath, cfg.GetProjectMainContainer())
			} else {
				fmt.Printf("Completed copying %s from container %s to host\n", syncPath, cfg.GetProjectMainContainer())
			}

			if direction == syncCopyTo {

				user := "root"

				if cfg.GetMainContainerUser() != "" {
					user = cfg.GetMainContainerUser()
				}

				p := strings.TrimRight(cfg.GetProjectDockerPath(), string(os.PathSeparator)) + syncPath
				args = []string{"exec", "-it", cfg.GetProjectMainContainer(), "chown", "-R", user + ":" + user, p}

				eo := ExecOptions{
					command: "docker",
					args:    args,
					tty:     true,
					detach:  true,
				}

				return execCommand(eo, c.App)
			}

			return nil
		},
	}
}
