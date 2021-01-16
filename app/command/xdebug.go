package command

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

type xdebug struct {
	usage       map[string]string
	aliases     map[string]string
	description map[string]string
}

type xdebugProjectConfig interface {
	GetXDebugFpmIniPath() string
	GetXDebugCliIniPath() string
	GetXDebugConfigLocaton() string
	GetProjectMainContainer() string
	SaveContainerNameToProjectConfig(string) error
	SaveDockerCliXdebugIniFilePath(string) error
	SaveDockerFpmXdebugIniFilePath(string) error
	SaveXDebugConifgLocaton(string) error
	GetCommandInactveStatus(string) bool
}

type xdebugArgsProjectConfig interface {
	GetXDebugFpmIniPath() string
	GetXDebugCliIniPath() string
	GetXDebugConfigLocaton() string
	GetProjectMainContainer() string
}

func getXdebugArgs(cfg xdebugArgsProjectConfig, command string, currentPath string) []string {

	action := `s/^\;zend_extension/zend_extension/g`

	if strings.Contains(command, "disable") {
		action = `s/^zend_extension/\;zend_extension/g`
	}

	xdebugFileConfigPath := cfg.GetXDebugFpmIniPath()

	if command == "xdebug:cli:enable" || command == "xdebug:cli:disable" {
		xdebugFileConfigPath = cfg.GetXDebugCliIniPath()
	}

	args := []string{}

	if cfg.GetXDebugConfigLocaton() == "local" {
		args = []string{"sed", "-i", "-e", action, strings.TrimRight(currentPath, string(os.PathSeparator)) + string(os.PathSeparator) + strings.Trim(xdebugFileConfigPath, string(os.PathSeparator))}
	} else {
		args = []string{"docker", "exec", cfg.GetProjectMainContainer(), "sed", "-i", "-e", action, xdebugFileConfigPath}
	}

	return args
}

type xDebugCommandDialog interface {
	SetMainContaner([]string) (int, string, error)
	DockerCliXdebugIniFilePath(string) (string, error)
	DockerFpmXdebugIniFilePath(string) (string, error)
	XDebugConfigLocation() (int, string, error)
}

type xDebugOptions interface {
	GetExecCommand() func(string, []string, *cli.App) error
	GetInitFuntion() func(bool) string
	GetContainerList() ([]string, error)
}

//XDebugCommand enable/disable xDebug
func XDebugCommand(xdebugAction string, cfg xdebugProjectConfig, d xDebugCommandDialog, options xDebugOptions) *cli.Command {
	execCommand := options.GetExecCommand()
	initf := options.GetInitFuntion()

	x := &xdebug{
		usage: map[string]string{
			"xdebug:fpm:enable":  "Enable fpm xdebug",
			"xdebug:fpm:disable": "Disable fpm xdebug",
			"xdebug:cli:enable":  "Enable cli xdebug",
			"xdebug:cli:disable": "Disable cli xdebug",
		},
		aliases: map[string]string{
			"xdebug:fpm:enable":  "xe",
			"xdebug:fpm:disable": "xd",
			"xdebug:cli:enable":  "xce",
			"xdebug:cli:disable": "xcd",
		},
		description: map[string]string{
			"xdebug:fpm:enable":  "",
			"xdebug:fpm:disable": "",
			"xdebug:cli:enable":  "",
			"xdebug:cli:disable": "",
		},
	}

	return &cli.Command{
		Name:            xdebugAction,
		Aliases:         []string{x.aliases[xdebugAction]},
		Usage:           x.usage[xdebugAction],
		Description:     x.description[xdebugAction],
		SkipFlagParsing: false,
		Hidden:          cfg.GetCommandInactveStatus("xdebug"),
		Action: func(c *cli.Context) (err error) {
			currentPath := initf(true)

			var cl []string

			if cl, err = options.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			if err = defineCliXdebugIniFilePath(cfg, d, "/etc/php/7.0/cli/conf.d/xdebug.ini"); err != nil {
				return err
			}

			if err = defineFpmXdebugIniFilePath(cfg, d, "/etc/php/7.0/fpm/conf.d/xdebug.ini"); err != nil {
				return err
			}

			if err = defineXdebugIniFileLocation(cfg, d); err != nil {
				return err
			}

			args := getXdebugArgs(cfg, xdebugAction, currentPath)

			if err = execCommand(args[0], args[1:], c.App); err != nil {
				return err
			}

			if xdebugAction == "xdebug:fpm:enable" {
				fmt.Printf("Fpm Xdebug enabled\n")
			} else if xdebugAction == "xdebug:fpm:disable" {
				fmt.Printf("Fpm Xdebug disabled \n")
			} else if xdebugAction == "xdebug:cli:enable" {
				fmt.Printf("Fpm Xdebug enabled \n")
			} else if xdebugAction == "xdebug:cli:disable" {
				fmt.Printf("Fpm Xdebug disabled \n")
			}

			return restartMainContainer(cfg)
		},
	}
}

type defineCliXdebugIniFilePathProjectConfig interface {
	SaveDockerCliXdebugIniFilePath(string) error
	GetXDebugCliIniPath() string
}

type defineCliXdebugIniFilePathDialog interface {
	DockerCliXdebugIniFilePath(string) (string, error)
}

func defineCliXdebugIniFilePath(cfg defineCliXdebugIniFilePathProjectConfig, d defineCliXdebugIniFilePathDialog, defaultPath string) (err error) {
	if cfg.GetXDebugCliIniPath() == "" {
		var path string
		if path, err = d.DockerCliXdebugIniFilePath(defaultPath); err != nil {
			return err
		}

		if path == "" {
			return errors.New("Cli Xdebug ini file path is empty")
		}

		return cfg.SaveDockerCliXdebugIniFilePath(path)
	}

	return nil
}

type defineFpmXdebugIniFilePathProjectConfig interface {
	SaveDockerFpmXdebugIniFilePath(string) error
	GetXDebugFpmIniPath() string
}

type defineFpmXdebugIniFilePathDialog interface {
	DockerFpmXdebugIniFilePath(string) (string, error)
}

func defineFpmXdebugIniFilePath(cfg defineFpmXdebugIniFilePathProjectConfig, d defineFpmXdebugIniFilePathDialog, defaultPath string) (err error) {
	if cfg.GetXDebugFpmIniPath() == "" {
		var path string
		if path, err = d.DockerFpmXdebugIniFilePath(defaultPath); err != nil {
			return err
		}

		if path == "" {
			return errors.New("Fpm Xdebug ini file path is empty")
		}

		return cfg.SaveDockerFpmXdebugIniFilePath(path)
	}

	return nil
}

type defineXdebugIniFileLocationProjectConfig interface {
	SaveXDebugConifgLocaton(string) error
	GetXDebugConfigLocaton() string
}

type defineXdebugIniFileLocationDialog interface {
	XDebugConfigLocation() (int, string, error)
}

func defineXdebugIniFileLocation(cfg defineXdebugIniFileLocationProjectConfig, d defineXdebugIniFileLocationDialog) (err error) {
	if cfg.GetXDebugConfigLocaton() == "" {
		var path string

		if _, path, err = d.XDebugConfigLocation(); err != nil {
			return err
		}

		if path == "" {
			return errors.New("Xdebug config file locaton cannot be empty")
		}

		return cfg.SaveXDebugConifgLocaton(path)
	}

	return nil
}
