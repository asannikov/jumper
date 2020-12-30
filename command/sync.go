package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

type sync struct {
	usage         map[string]string
	aliases       map[string]string
	description   map[string]string
	initFunction  func()
	containerList []string
	ctype         string
	command       string
}

//SyncCommand does the syncronization between container and project
func SyncCommand(direction string, initf func(bool), dockerStatus bool, cfg projectConfig, d dialog, clist containerlist) *cli.Command {

	s := &sync{
		usage: map[string]string{
			"copyto":   "Runs composer: {docker exec -it phpContainer composer} [custom parameters]",
			"copyfrom": "Runs composer with no memory constraint: {docker exec -i phpContainer /usr/bin/php -d memory_limit=-1 /usr/local/bin/composer} [custom parameters]",
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
			initf(true)

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

			pathFrom := "./"
			args := []string{"cp", pathFrom, cfg.GetProjectMainContainer() + ":" + cfg.GetProjectDockerPath()}

			cmd := exec.Command("docker", args...)

			fmt.Printf("\ncommand: %s\n\n", "docker "+strings.Join(args, " "))
			os.Exit(1)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err = cmd.Run(); err != nil {
				return err
			}

			fmt.Printf("Completed copying %s files from host to container\n", "all")

			return nil
		},
	}
}

// 1. define docker project root dialog
// 2.

/*
#!/bin/bash
[ -z "$1" ] && echo "Please specify a directory or file to copy to container (ex. vendor, --all)" && exit

REAL_SRC=$(cd -P "src" && pwd)
if [ "$1" == "--all" ]; then
  docker cp $REAL_SRC/./ $(docker-compose ps -q phpfpm|awk '{print $1}'):/var/www/html/
  echo "Completed copying all files from host to container"
  bin/fixowns
  bin/fixperms
else
  if [ -f "$REAL_SRC/$1" ]; then
    docker cp $REAL_SRC/$1 $(docker-compose ps -q phpfpm|awk '{print $1}'):/var/www/html/$1
  else
    docker cp $REAL_SRC/$1 $(docker-compose ps -q phpfpm|awk '{print $1}'):/var/www/html/$(dirname $1)
  fi
  echo "Completed copying $1 from host to container"
  bin/fixowns $1
  bin/fixperms $1
fi
*/
