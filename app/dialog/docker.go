package dialog

import (
	"fmt"
	"runtime"

	"github.com/manifoldco/promptui"
)

// DockerService call the request dialog to define docker service
func (d *Dialog) DockerService() (string, error) {
	return d.setDockerService()
}

// StartDocker call the request dialog to start docker
func (d *Dialog) StartDocker() (string, error) {
	return d.setStartDocker()
}

// StartCommand sets start docker command
func (d *Dialog) StartCommand() (string, error) {
	return d.setStartCommand()
}

func setStartCommand() (string, error) {
	validate := func(c string) error {
		if c == "" {
			return fmt.Errorf("Command name cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Set start command",
		Validate: validate,
		Default:  "docker-compose -f docker-compose.yml up --force-recreate -d --remove-orphans",
	}

	return prompt.Run()
}

func dockerService() (string, error) {
	validate := func(c string) error {
		if c == "" {
			return fmt.Errorf("Command name cannot be empty")
		}
		return nil
	}

	var defaultCommand string

	if runtime.GOOS == "linux" {
		defaultCommand = "service docker start"
	} else if runtime.GOOS == "darwin" {
		defaultCommand = "open --hide -a Docker"
	} else if runtime.GOOS == "windows" {
		defaultCommand = `C:\Program Files\docker\dockerd.exe`
	}

	prompt := promptui.Prompt{
		Label:    "Set docker service command",
		Validate: validate,
		Default:  defaultCommand,
	}

	return prompt.Run()
}

func startDocker() (string, error) {
	prompt := promptui.Prompt{
		Label:     "Start Docker",
		IsConfirm: true,
		Default:   "y",
	}

	return prompt.Run()
}
