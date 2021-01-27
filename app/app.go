package app

import (
	"fmt"
	"io"
	"os"

	"github.com/asannikov/jumper/app/lib"

	"github.com/asannikov/jumper/app/config"
	"github.com/asannikov/jumper/app/dialog"

	"github.com/urfave/cli/v2"
)

const confgFile = "jumper.json"

// CommandList defines type for command list in main
type CommandList = func(*config.Config, *dialog.Dialog, bool, func(string, string) (string, error), func(bool) string) []*cli.Command

// JumperApp initializate app
func JumperApp(cli *cli.App) {
	// Dialogs
	DLG := dialog.InitDialogFunctions()

	cfg := &config.Config{
		ProjectFile: confgFile,
	}
	cfg.Init()

	fs := &FileSystem{}
	cfg.SetFileSystem(fs)

	// Loading only global config
	loadGlobalConfig(cfg, fs)

	// Loading project config if exists
	loadProjectConfig(cfg, fs)

	// Define docker command
	defineDockerCommand(cfg, &DLG)

	initf := func(seekProject bool) string {
		if err := seekPath(cfg, &DLG, fs, seekProject); err != nil {
			fmt.Println(err)
			os.Exit(1)
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

type ioCli struct {
	reader    io.Reader
	writer    io.Writer
	errWriter io.Writer
}

func (ic *ioCli) GetReader() io.Reader {
	return ic.reader
}

func (ic *ioCli) GetWriter() io.Writer {
	return ic.writer
}

func (ic *ioCli) GetErrWriter() io.Writer {
	return ic.errWriter
}
