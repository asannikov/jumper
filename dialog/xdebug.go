package dialog

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// DockerCliXdebugIniFilePath gets the path to xdebug cli ini file
func (d *Dialog) DockerCliXdebugIniFilePath(defaulPath string) (string, error) {
	return d.setDockerCliXdebugIniFilePath(defaulPath)
}

// DockerFpmXdebugIniFilePath gets the path to xdebug fpm ini file
func (d *Dialog) DockerFpmXdebugIniFilePath(defaulPath string) (string, error) {
	return d.setDockerFmpXdebugIniFilePath(defaulPath)
}

func dockerCliXdebugIniFilePath(path string) (string, error) {
	validate := func(p string) error {
		if p == "" {
			return fmt.Errorf("Xdebug cli file path cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Add path to xdebug cli ini config",
		Validate: validate,
		Default:  path,
	}

	return prompt.Run()
}

func dockerFmpXdebugIniFilePath(path string) (string, error) {
	validate := func(p string) error {
		if p == "" {
			return fmt.Errorf("Xdebug fpm file path cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Add path to xdebug fpm ini config",
		Validate: validate,
		Default:  path,
	}

	return prompt.Run()
}
