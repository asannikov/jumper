package main

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

type dialogInternal struct {
	setPhpContaner func() (int, string, error)
	setProjectPath func() (string, error)
	setProjectName func() (string, error)
	selectProject  func([]string) (int, string, error)
	addProjectPath func(string) (string, error)
	addProjectName func() (string, error)
}

func initDialogFunctions() dialogInternal {
	return dialogInternal{
		setPhpContaner: selectPhpContainer,
		setProjectPath: selectProjectPath,
		selectProject:  selectProject,
		addProjectPath: addProjectPath,
		addProjectName: addProjectName,
	}
}

func selectPhpContainer() (int, string, error) {
	containers := getContainerList()

	prompt := promptui.Select{
		Label: "Select php container",
		Items: containers,
	}

	return prompt.Run()
}

// get project path or create
func selectProjectPath() (string, error) {
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
