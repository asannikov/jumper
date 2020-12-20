package main

import (
	"fmt"
	"log"
	"jumper/config"
	"jumper/container"
	"jumper/dialog"
	"os"

	"github.com/urfave/cli/v2"
	// imports as package "cli"
)

const confgFile = "jumper.json"

func main() {

	// Dialogs
	DLG := dialog.InitDialogFunctions()

	cfg := &config.Config{
		ProjectFile: confgFile,
	}
	cfg.Init()

	fs := &FileSystem{}
	cfg.SetFileSystem(fs)

	initf := func() {
		if err := seekPath(cfg, &DLG, fs); err != nil {
			log.Fatal(err)
		}

		currentDir, _ := fs.GetWd()

		fmt.Printf("\nchanged user location to directory: %s\n\n", currentDir)
	}

	app := &cli.App{
		Commands: getCommandList(cfg, &DLG, initf),
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func getContainerList() []string {
	return container.GetContanerList()
}
