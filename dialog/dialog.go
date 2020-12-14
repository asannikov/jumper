package dialog

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// Dialog contains methods for the iteraction with promptui
type Dialog struct {
	setMainContaner func([]string) (int, string, error)
	setStartCommand func() (string, error)

	//SetProjectPath func() (string, error)
	//SetProjectName func() (string, error)

	// Project management
	SelectProject  func([]string) (int, string, error)
	AddProjectPath func(string) (string, error)
	AddProjectName func() (string, error)
}

// SetStartCommand sets main container name
func (d *Dialog) SetStartCommand() (string, error) {
	return d.setStartCommand()
}

// SetMainContaner sets main container name
func (d *Dialog) SetMainContaner(cl []string) (int, string, error) {
	return d.setMainContaner(cl)
}

// InitDialogFunctions initiate all methods
func InitDialogFunctions() Dialog {
	return Dialog{
		//SetPhpContaner: selectPhpContainer,
		//SetProjectPath: selectProjectPath,
		SelectProject:  selectProject,
		AddProjectPath: addProjectPath,
		AddProjectName: addProjectName,

		setMainContaner: setMainContaner,
		setStartCommand: setStartCommand,
	}
}

type projectConfig interface {
	GetProjectName() string
	GetProjectPath() string
	SetProjectName(string)
	SetProjectPath(string)
}

// CallAddProjectDialog calls project manager
func (d *Dialog) CallAddProjectDialog(pc projectConfig) error {
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

// select project path from the list
func selectProject(projects []string) (int, string, error) {
	prompt := promptui.SelectWithAdd{
		Label:    "Select project from the list",
		Items:    projects,
		AddLabel: "Add new project",
	}

	return prompt.Run()
}

func addProjectPath(path string) (string, error) {
	validate := func(p string) error {
		if p == "" {
			return fmt.Errorf("Project name cannot be empty")
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

func setMainContaner(containers []string) (int, string, error) {
	prompt := promptui.Select{
		Label: "Select main container",
		Items: containers,
	}

	return prompt.Run()
}

func setStartCommand() (string, error) {
	validate := func(c string) error {
		if c == "" {
			return fmt.Errorf("Command name cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Set start start",
		Validate: validate,
		Default:  "docker-compose -f docker-compose.yml up --force-recreate -d --remove-orphans $1",
	}

	return prompt.Run()
}
