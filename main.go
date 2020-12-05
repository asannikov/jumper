package main

import (
	"log"
	"mgt/config"
	"mgt/container"
	"os"
	"os/user"

	"github.com/urfave/cli/v2"
	// imports as package "cli"
)

const confgFile = "mgt.json"

func main() {
	// Dialogs
	DLG := initDialogFunctions()

	userDir, err := getUserDirectory()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.MgtConfig{
		FileName: confgFile,
		FileUser: userDir + string(os.PathSeparator) + ".mgt.json",
		FileSystem: config.FileSystem{
			FileExists:     fileExist,
			DirExists:      dirExists,
			ReadConfigFile: readConfigFile,
			SaveConfigFile: saveConfigFile,
		},
	}

	app := &cli.App{
		Commands: getCommandList(&cfg, &DLG),
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getContainerList() []string {
	return container.GetContanerList()
}

func getUserDirectory() (string, error) {
	usr, err := user.Current()

	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
