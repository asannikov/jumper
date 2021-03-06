package app

import (
	"errors"
	"os"

	"github.com/asannikov/jumper/app/config"
	"github.com/asannikov/jumper/app/dialog"
)

type projectConfig struct {
	projectpath        string
	projectname        string
	selectProjectIsNew bool
}

func (pc *projectConfig) GetProjectName() string {
	return pc.projectname
}

func (pc *projectConfig) GetProjectPath() string {
	return pc.projectpath
}

func (pc *projectConfig) SetProjectName(v string) {
	pc.projectname = v
}

func (pc *projectConfig) SetProjectPath(v string) {
	pc.projectpath = v
}

type definePathsFileSystem interface {
	GetUserDirectory() (string, error)
}

type definePathsConfig interface {
	SetUserFile(string)
}

func definePaths(cfg loadGlobalCfg, fs definePathsFileSystem) (err error) {
	userDir, err := fs.GetUserDirectory()
	cfg.SetUserFile(userDir + string(os.PathSeparator) + ".jumper.json")
	return err
}

type runDialogConfig interface {
	FindProjectPathInJSON(config.ProjectSettings)
	ProjectConfgFileFound() bool
}

type runDialogDialog interface {
	SelectProject([]string) (int, string, error)
	CallAddProjectDialog(dialog.ProjectConfig) error
}

func runDialog(pc *projectConfig, cfg runDialogConfig, DLG runDialogDialog, pl []string, currentDir string) (err error) {
	if cfg.ProjectConfgFileFound() == false && len(pl) > 0 {
		var index int
		var projectName string

		if index, projectName, err = DLG.SelectProject(pl); err != nil {
			return err
		}

		pc.SetProjectName(projectName)

		if index == -1 { // add new project
			pc.SetProjectPath(currentDir)
			err = DLG.CallAddProjectDialog(pc)
			pc.selectProjectIsNew = true
		} else {
			cfg.FindProjectPathInJSON(pc)
		}

		if err != nil {
			return err
		}
	}

	if pc.GetProjectPath() == "" {
		pc.selectProjectIsNew = true
		pc.SetProjectPath(currentDir)
		err = DLG.CallAddProjectDialog(pc)
	}

	return err
}

type seekPathFileSystem interface {
	FileExists(string) (bool, error)
	DirExists(string) (bool, error)
	SaveConfigFile(interface{}, string) error
	ReadConfigFile(string, interface{}) error
	GoToProjectPath(string) error
	GetUserDirectory() (string, error)
	GetWd() (string, error)
}

type defineDockerCommandConfig interface {
	SetDockerCommand(string) error
	GetDockerCommand() string
}

type defineDockerCommandDialog interface {
	DockerService() (string, error)
}

func defineDockerCommand(cfg defineDockerCommandConfig, DLG defineDockerCommandDialog) (err error) {
	command := cfg.GetDockerCommand()
	if command == "" {
		if command, err = DLG.DockerService(); err != nil {
			return err
		}

		cfg.SetDockerCommand(command)
	}
	return nil
}

type loadGlobalCfg interface {
	LoadConfig(bool) error
	SetUserFile(string)
}

type loadGlobalSystem interface {
	GetUserDirectory() (string, error)
}

func loadGlobalConfig(cfg loadGlobalCfg, fs loadGlobalSystem) (err error) {
	if err = definePaths(cfg, fs); err != nil {
		return err
	}

	if err = cfg.LoadConfig(false); err != nil {
		return err
	}

	return nil
}

type loadProjectCfg interface {
	LoadProjectConfig() (bool, error)
	GetProjectName() string
	SetProjectPath(string)
}

type loadProjectFs interface {
	GetWd() (string, error)
}

func loadProjectConfig(cfg loadProjectCfg, fs loadProjectFs) (err error) {
	var currentDir string
	var status bool

	if currentDir, err = fs.GetWd(); err != nil {
		return err
	}

	status, _ = cfg.LoadProjectConfig()

	if status && cfg.GetProjectName() != "" {
		cfg.SetProjectPath(currentDir)
	}

	return nil
}

type seekPathDialog interface {
	CallAddProjectDialog(dialog.ProjectConfig) error
	SelectProject([]string) (int, string, error)
}

func seekPath(cfg *config.Config, DLG seekPathDialog, fs seekPathFileSystem, seekProject bool) error {
	var currentDir string
	var err error

	if !seekProject {
		return nil
	}

	if err = loadGlobalConfig(cfg, fs); err != nil {
		return err
	}

	if currentDir, err = fs.GetWd(); err != nil {
		return err
	}

	pl := cfg.GetProjectNameList()

	pc := &projectConfig{
		projectpath: cfg.GetProjectPath(),
		projectname: cfg.GetProjectName(),
	}

	if err = runDialog(pc, cfg, DLG, pl, currentDir); err != nil {
		return err
	}

	cfg.SetProjectName(pc.GetProjectName())
	cfg.SetProjectPath(pc.GetProjectPath())

	if pc.selectProjectIsNew == true {
		err = cfg.AddProjectConfigFile()
	}

	if err != nil {
		return err
	}

	if pc.GetProjectPath() == "" {
		return errors.New("No project found")
	}

	if err = fs.GoToProjectPath(pc.GetProjectPath()); err != nil {
		return err
	}

	return cfg.LoadConfig(seekProject)
}
