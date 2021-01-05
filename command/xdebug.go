package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

type xdebug struct {
	usage       map[string]string
	aliases     map[string]string
	description map[string]string
}

func getXdebugArgs(cfg projectConfig, command string, currentPath string) []string {

	action := `s/^\;zend_extension/zend_extension/g`

	if strings.Contains(command, "disable") {
		action = `s/^zend_extension/\;zend_extension/g`
	}

	xdebugFileConfigPath := cfg.GetXDebugFpmIniPath()

	if command == "xdebug:cli:enable" || command == "xdebug:cli:disable" {
		xdebugFileConfigPath = cfg.GetXDebugCliIniPath()
	}

	args := []string{}

	if cfg.GetXDebugConifgLocaton() == "local" {
		args = []string{"sed", "-i", "-e", action, currentPath + string(os.PathSeparator) + strings.Trim(xdebugFileConfigPath, string(os.PathSeparator))}
	} else {
		args = []string{"docker", "exec", cfg.GetProjectMainContainer(), "sed", "-i", "-e", action, xdebugFileConfigPath}
	}

	return args
}

//XDebugCommand enable/disable xDebug
func XDebugCommand(xdebugAction string, initf func(bool) string, dockerStatus bool, cfg projectConfig, d dialog, clist containerlist) *cli.Command {

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
		Action: func(c *cli.Context) (err error) {
			currentPath := initf(true)

			var cl []string

			if cl, err = clist.GetContainerList(); err != nil {
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

			cmd := exec.Command(args[0], args[1:]...)

			fmt.Printf("\ncommand: %s\n\n", args[0]+" "+strings.Join(args[1:], " "))

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err = cmd.Run(); err != nil {
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

			return nil
		},
	}
}

/*
docker-compose exec phpfpm sed -i -e 's/^zend_extension/\;zend_extension/g' /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini
docker-compose exec phpfpm  sed -i -e 's/^\;zend_extension/zend_extension/g' /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini

#!/bin/bash
if [ "$1" == "disable" ]; then
  bin/cli sed -i -e 's/^zend_extension/\;zend_extension/g' /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini
  sleep 1
  bin/restart phpfpm
  echo "Xdebug has been disabled."
elif [ "$1" == "enable" ]; then
  bin/cli sed -i -e 's/^\;zend_extension/zend_extension/g' /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini
  sleep 1
  bin/restart phpfpm
  echo "Xdebug has been enabled."
else
  echo "Please specify either 'enable' or 'disable' as an argument"
fi

*/
