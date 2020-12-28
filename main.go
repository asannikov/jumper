package main

import (
	"fmt"
	"jumper/config"
	"jumper/dialog"
	"jumper/lib"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	// imports as package "cli"
)

const confgFile = "jumper.json"
const version = "1.2.0"

func main() {

	// Dialogs
	DLG := dialog.InitDialogFunctions()

	cfg := &config.Config{
		ProjectFile: confgFile,
	}
	cfg.Init()

	fs := &FileSystem{}
	cfg.SetFileSystem(fs)

	// Loading only global config
	loadGlobalConfig(cfg, &DLG, fs)

	// Define docker command
	defineDockerCommand(cfg, &DLG)

	initf := func(seekProject bool) {
		if err := seekPath(cfg, &DLG, fs, seekProject); err != nil {
			log.Fatal(err)
		}

		if seekProject == true {
			currentDir, _ := fs.GetWd()
			fmt.Printf("\nchanged user location to directory: %s\n\n", currentDir)
		}
	}

	app := &cli.App{
		Name:                 "Jumper - the tool for quick docker project management in cli",
		Usage:                "jumper [command] [parameters]",
		Description:          "Create project config using Jumper and work with your docker project in cli without routine",
		EnableBashCompletion: true,
		Commands:             getCommandList(cfg, &DLG, initf),
		Version:              version,
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Anton Sannikov",
				Email: "",
			},
		},
		Copyright: lib.GetCopyrightText(cfg),
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
