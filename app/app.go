package app

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/asannikov/jumper/app/lib"

	"github.com/asannikov/jumper/app/config"
	"github.com/asannikov/jumper/app/dialog"

	"github.com/urfave/cli/v2"
)

const confgFile = "jumper.json"

// CommandList defines type for command list in main
type CommandList = func(*config.Config, *dialog.Dialog, bool, func(string, string) (string, error), func(bool) string) []*cli.Command

// InitApp initializate app
func InitApp(cli *cli.App) {
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

	// Loading project config if exists
	loadProjectConfig(cfg, fs)

	// Define docker command
	defineDockerCommand(cfg, &DLG)

	initf := func(seekProject bool) string {
		if err := seekPath(cfg, &DLG, fs, seekProject); err != nil {
			log.Fatal(err)
		}

		if seekProject == true {
			currentDir, _ := fs.GetWd()
			fmt.Printf("\nchanged user location to directory: %s\n\n", currentDir)
			return currentDir
		}

		return ""
	}

	cli.Copyright = lib.GetCopyrightText(cfg)
	cli.Commands = commandList(cfg, &DLG, initf)
}

func execCommand(command string, args []string, app *cli.App) error {
	cmd := exec.Command(command, args...)

	cmd.Stdin = app.Reader
	cmd.Stdout = app.Writer
	cmd.Stderr = app.ErrWriter

	fmt.Printf("\ncommand: %s\n\n", command+" "+strings.Join(args, " "))

	return cmd.Run()
}
