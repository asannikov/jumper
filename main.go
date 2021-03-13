package main

import (
	"fmt"
	"os"

	"github.com/asannikov/jumper/app"

	"github.com/urfave/cli/v2"
)

const version = "1.8.8"

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

	app.JumperApp(c)

	err := c.Run(os.Args)

	if err != nil {
		fmt.Println(err)
	}
}
