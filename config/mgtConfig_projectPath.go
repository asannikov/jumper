package config

import (
	"fmt"
	"os"
	"strings"
)

type createProjectDialog struct {
	projectName string
	projectPath string
}

// GetProjectPath gets project path
func (c *MgtConfig) GetProjectPath(gf func() error) (string, error) {
	err := gf()

	if err != nil {
		return "", err
	}

	err = c.handleConfig()

	if err != nil {
		return "", err
	}
	return c.projectSettings.Path, nil
}

func projectListManagementDialog(apn func() (string, error), app func(string) (string, error), d *createProjectDialog) error {
	if d.projectName == "" {
		pn, err := apn() // add project name
		if err != nil {
			return err
		}
		d.projectName = pn
	}

	pp, err := app(d.projectPath) // add project path
	if err != nil {
		return err
	}
	d.projectPath = pp

	return nil
}

// EvaluateProjectPath helps to find/define project path
func (c *MgtConfig) EvaluateProjectPath(
	sp func([]string) (int, string, error),
	apn func() (string, error),
	app func(string) (string, error)) error {

	gc := GlobalConfig{}

	var t bool
	var err error
	createNewProject := false

	if t, err = c.FileSystem.FileExists(c.FileUser); err == nil && t == true {
		err = c.FileSystem.ReadConfigFile(c.FileUser, &gc)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	pl := gc.GetProjectNameList()

	cpd := createProjectDialog{}

	if len(pl) > 0 {
		index, projectName, errd := sp(pl) // select project in json or create new one

		if errd != nil {
			return errd
		}

		var foundPath string

		if index == -1 { // add new project
			currentDir, erros := os.Getwd()

			if erros != nil {
				return erros
			}
			cpd.projectPath = currentDir
			cpd.projectName = projectName
			err = projectListManagementDialog(apn, app, &cpd)
			foundPath = cpd.projectPath
			createNewProject = true
		} else {
			err = gc.FindProjectPathInJSON(func(p ProjectSettings) (bool, error) {
				if p.Name == projectName {
					if t, err = c.FileSystem.DirExists(p.Path); err == nil && t == true {
						foundPath = p.Path
						cpd.projectPath = p.Path
						cpd.projectName = p.Name
						return true, nil
					} else if err != nil {
						return false, err
					}
				}
				return false, nil
			})
		}

		if err != nil {
			return err
		}

		if foundPath == "" {
			return fmt.Errorf("Cannot use project in the root path")
		}

		foundPath = strings.TrimRight(foundPath, string(os.PathSeparator))
		if err := os.Chdir(foundPath); err != nil {
			return err
		}

		currentDir, err := os.Getwd()

		if err != nil {
			return err
		}

		if currentDir != foundPath {
			return fmt.Errorf("Expected path %s, the current one %s", currentDir, foundPath)
		}
	} else {
		if err := projectListManagementDialog(apn, app, &cpd); err != nil {
			return err
		}
	}

	p := ProjectSettings{
		Path: cpd.projectPath,
		Name: cpd.projectName,
	}

	if createNewProject == true {
		gc.AddNewProject(p)
	}

	if err := c.FileSystem.SaveConfigFile(gc, c.FileUser); err != nil {
		return err
	}

	c.globalSettings = gc
	c.projectSettings = p

	return nil
}
