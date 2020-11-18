package main

import (
	"log"
	"mgt/container"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2" // imports as package "cli"
)

const confgFile = "mgt.json"

// Configuration contains file config
type Configuration struct {
	PhpContainer string
}

type dialogInternal struct {
	setPhpContaner func() (int, string, error)
}

func main() {

	cfg := MgtConfig{
		FileName: confgFile,
	}

	app := &cli.App{
		Commands: getCommandList(&cfg),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getContainerList() []string {
	return container.GetContanerList()
}

func selectPhpContainer() (int, string, error) {
	containers := getContainerList()

	prompt := promptui.Select{
		Label: "Select php container",
		Items: containers,
	}

	return prompt.Run()
}
