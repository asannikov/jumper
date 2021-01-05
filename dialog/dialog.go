package dialog

import (
	"github.com/manifoldco/promptui"
)

// Dialog contains methods for the iteraction with promptui
type Dialog struct {
	setMainContaner               func([]string) (int, string, error)
	setStartCommand               func() (string, error)
	setStartDocker                func() (string, error)
	setDockerService              func() (string, error)
	setDockerProjectPath          func(string) (string, error)
	setDockerCliXdebugIniFilePath func(string) (string, error)
	setDockerFmpXdebugIniFilePath func(string) (string, error)
	setXdebugFileConfigLocation   func() (int, string, error)

	// Project management
	SelectProject  func([]string) (int, string, error)
	AddProjectPath func(string) (string, error)
	AddProjectName func() (string, error)
}

// InitDialogFunctions initiate all methods
func InitDialogFunctions() Dialog {
	return Dialog{
		SelectProject:  selectProject,
		AddProjectPath: addProjectPath,
		AddProjectName: addProjectName,

		setMainContaner:               setMainContaner,
		setStartCommand:               setStartCommand,
		setStartDocker:                startDocker,
		setDockerService:              dockerService,
		setDockerProjectPath:          dockerProjectPath,
		setDockerCliXdebugIniFilePath: dockerCliXdebugIniFilePath,
		setDockerFmpXdebugIniFilePath: dockerFmpXdebugIniFilePath,
		setXdebugFileConfigLocation:   xdebugFileConfigLocation,
	}
}

type projectConfig interface {
	GetProjectName() string
	GetProjectPath() string
	SetProjectName(string)
	SetProjectPath(string)
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
