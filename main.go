package main

import (
	"log"
	"os"

	"jumper/app" // github.com/asannikov/

	"github.com/urfave/cli/v2"
)

const version = "1.5.5"

func main() {

	c := &cli.App{
		Name:                 "Jumper - the tool for quick docker project management in cli",
		Usage:                "jumper [command] [parameters]",
		Description:          "Create project config using Jumper and work with your docker project in cli without routine",
		EnableBashCompletion: true,
		Version:              version,
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Anton Sannikov",
				Email: "",
			},
		},
	}

	app.InitApp(c)

	err := c.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
