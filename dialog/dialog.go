package dialog

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// Dialog contains methods for the iteraction with promptui
type Dialog struct {
	//SetPhpContaner func([]string) (int, string, error)
	//SetProjectPath func() (string, error)
	//SetProjectName func() (string, error)
	SelectProject  func([]string) (int, string, error)
	AddProjectPath func(string) (string, error)
	AddProjectName func() (string, error)
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

// InitDialogFunctions initiate all methods
func InitDialogFunctions() Dialog {
	return Dialog{
		//SetPhpContaner: selectPhpContainer,
		//SetProjectPath: selectProjectPath,
		SelectProject:  selectProject,
		AddProjectPath: addProjectPath,
		AddProjectName: addProjectName,
	}
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

/*func selectPhpContainer(containers []string) (int, string, error) {
	prompt := promptui.Select{
		Label: "Select php container",
		Items: containers,
	}

	return prompt.Run()
}*/

// get project path or create
/*func selectProjectPath() (string, error) {
	//return "Volumes/LS/AV/golang/backendarmy-blog/mgt", nil
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	validate := func(p string) error {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			return err
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Path to project",
		Validate: validate,
		Default:  path,
	}

	return prompt.Run()
}

// set project name
func setProjectName() (string, error) {
	validate := func(p string) error {
		if p == "" {
			return fmt.Errorf("Project name cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Project name",
		Validate: validate,
	}

	return prompt.Run()
}*/
