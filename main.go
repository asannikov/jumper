package main

import (
	"fmt"
	"log"
	"mgt/config"
	"mgt/container"
	"mgt/dialog"
	"os"

	"github.com/urfave/cli/v2"
	// imports as package "cli"
)

const confgFile = "mgt.json"

func main() {

	// Dialogs
	DLG := dialog.InitDialogFunctions()

	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	fs := &FileSystem{}
	cfg.SetFileSystem(fs)

	if err := seekPath(cfg, &DLG, fs); err != nil {
		log.Fatal(err)
	}

	currentDir, _ := fs.GetWd()

	fmt.Printf("changed user location to directory: %s\n\n", currentDir)

	app := &cli.App{
		Commands: getCommandList(cfg, &DLG),
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func getContainerList() []string {
	return container.GetContanerList()
}
