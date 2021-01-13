package dialog

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// DockerProjectPath gets the path to container
func (d *Dialog) DockerProjectPath(defaulPath string) (string, error) {
	return d.setDockerProjectPath(defaulPath)
}

func dockerProjectPath(path string) (string, error) {
	validate := func(p string) error {
		if p == "" {
			return fmt.Errorf("Project path cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Add project path in docker container",
		Validate: validate,
		Default:  path,
	}

	return prompt.Run()
}
