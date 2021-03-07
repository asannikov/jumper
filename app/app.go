package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/asannikov/jumper/app/bash"
	"github.com/asannikov/jumper/app/command"
	"github.com/asannikov/jumper/app/docker"
	"github.com/asannikov/jumper/app/lib"
	"github.com/docker/docker/api/types"

	"github.com/asannikov/jumper/app/config"
	"github.com/asannikov/jumper/app/dialog"

	"github.com/urfave/cli/v2"
)

const confgFile = "jumper.json"

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

	opt := getOptions(cfg, &DLG)
	opt.setInitFuntion(func(seekProject bool) string {
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
	})

	cli.Copyright = lib.GetCopyrightText(cfg)
	cli.Commands = commandList(cfg, &DLG, opt)
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

type getOptionsConfig interface {
	GetDockerCommand() string
	GetProjectMainContainer() string
}

func getOptions(c getOptionsConfig, d commandListDialog) *commandOptions {
	b := bash.Bash{}

	getCommandLocationF := b.GetCommandLocation()

	dck := docker.GetDockerInstance()

	dockerDialog := getDockerStartDialog()
	dockerDialog.setDialog(d)
	dockerDialog.setDocker(dck)
	dockerDialog.setDockerService(c.GetDockerCommand())

	dockerStatus := false

	if dockerAPIVersiongo, _ := dck.Stat(); dockerAPIVersiongo != "" {
		dockerStatus = true
		dck.InitClient()
	}

	opt := &commandOptions{}
	opt.setCommandLocation(getCommandLocationF)
	opt.setDockerStatus(dockerStatus)
	opt.setStopContainers(dck.StopContainers())
	opt.setExecCommand(func(eo command.ExecOptions, app *cli.App) error {

		cmd := exec.Command(eo.GetCommand(), eo.GetArgs()...)

		cmd.Stdin = app.Reader
		cmd.Stdout = app.Writer
		cmd.Stderr = app.ErrWriter

		fmt.Printf("\ncommand: %s\n\n", eo.GetCommand()+" "+strings.Join(eo.GetArgs(), " "))

		return cmd.Run()
	})
	opt.setDirExists(func(path string) (bool, error) {
		if info, err := os.Stat(path); err == nil {
			if info.IsDir() {
				return true, nil
			}
			return false, fmt.Errorf("Path %s is a file ", path)
		} else if os.IsNotExist(err) {
			// path does *not* exist
			return false, err
		} else {
			// Schrodinger: file may or may not exist. See err for details.
			// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
			return false, err
		}
	})
	opt.setMkdirAll(func(path string, fileMode os.FileMode) error {
		return os.MkdirAll(path, fileMode)
	})

	opt.setDockerDialog(dockerDialog)
	opt.setCopyTo(func(container string, sourcePath string, dstPath string) error {
		return dck.CopyTo(container, sourcePath, dstPath)
	})
	opt.setNativeExec(func(eo command.ExecOptions, app *cli.App) (err error) {
		ic := &ioCli{
			reader:    app.Reader,
			writer:    app.Writer,
			errWriter: app.ErrWriter,
		}

		cnf := types.ExecConfig{
			AttachStderr: true,
			AttachStdin:  true,
			AttachStdout: true,
			User:         eo.GetUser(),
			Tty:          eo.GetTty(),
			Cmd:          append([]string{eo.GetCommand()}, eo.GetArgs()...),
			WorkingDir:   eo.GetWorkingDir(),
		}

		status, err := dck.Exec(c.GetProjectMainContainer(), &cnf, ic)

		if err != nil {
			return err
		}

		if status > 0 {
			return errors.New("Error is occurred on exec function")
		}

		return nil
	})
	opt.setMagentoBin(command.CheckMagentoBin)
	opt.setXdebugStatus(command.CheckXdebugStatus)

	return opt
}
