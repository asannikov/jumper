package main

import (
	"fmt"
	"jumper/config"
	"jumper/container"
	"jumper/dialog"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	// imports as package "cli"
)

const confgFile = "jumper.json"
const version = "1.0.1"

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

	currentTime := time.Now()

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
		Copyright: fmt.Sprintf(`
MIT License

Copyright (c) %d Anton Sannikov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`, currentTime.Year()),
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func getContainerList() []string {
	return container.GetContanerList()
}
