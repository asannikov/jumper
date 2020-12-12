package main

import (
	"errors"
	"mgt/config"
	"mgt/dialog"
	"os"
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

func definePaths(cfg *config.Config, fs definePathsFileSystem) (err error) {
	userDir, err := fs.GetUserDirectory()
	cfg.UserFile = userDir + string(os.PathSeparator) + ".mgt.json"
	return err
}

type runDialogFileSystem interface {
	GoToProjectPath(string) error
}

func runDialog(pc *projectConfig, cfg *config.Config, DLG *dialog.Dialog, fs runDialogFileSystem, pl []string, currentDir string) (err error) {
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

func seekPath(cfg *config.Config, DLG *dialog.Dialog, fs seekPathFileSystem) error {
	var currentDir string
	var err error

	if err = definePaths(cfg, fs); err != nil {
		return err
	}

	if err = cfg.LoadConfig(); err != nil {
		return err
	}

	if currentDir, err = fs.GetWd(); err != nil {
		return err
	}

	pl := cfg.GetProjectNameList()

	pc := &projectConfig{}

	if err = runDialog(pc, cfg, DLG, fs, pl, currentDir); err != nil {
		return err
	}

	if pc.selectProjectIsNew == true {
		project := &config.ProjectConfig{
			Path: pc.GetProjectPath(),
			Name: pc.GetProjectName(),
		}
		err = cfg.AddProjectConfigFile(project)
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

	return nil
}
