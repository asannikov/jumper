package dialog

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// CallAddProjectDialog calls project manager
func (d *Dialog) CallAddProjectDialog(pc ProjectConfig) error {
	if pc.GetProjectName() == "" {
		pn, err := d.AddProjectName() // add project name
		if err != nil {
			return err
		}
		pc.SetProjectName(pn)
	}

	pp, err := d.AddProjectPath(pc.GetProjectPath()) // add project path
	if err != nil {
		return err
	}

	pc.SetProjectPath(pp)

	return nil
}

// AddProjectPath sets project path
func (d *Dialog) AddProjectPath(path string) (string, error) {
	return d.setAddProjectPath(path)
}

// AddProjectName sets project name
func (d *Dialog) AddProjectName() (string, error) {
	return d.setAddProjectName()
}

func addProjectPath(path string) (string, error) {
	validate := func(p string) error {
		if p == "" {
			return fmt.Errorf("Project path cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Add project path",
		Validate: validate,
		Default:  path,
	}

	return prompt.Run()
}

func addProjectName() (string, error) {
	validate := func(p string) error {
		if p == "" {
			return fmt.Errorf("Project name cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Add project name",
		Validate: validate,
	}

	return prompt.Run()
}
