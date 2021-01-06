package dialog

import "github.com/manifoldco/promptui"

// DockerShell defines shell type of docker main container
func (d *Dialog) DockerShell() (int, string, error) {
	return d.setDockerShell()
}

func dockerShell() (int, string, error) {
	prompt := promptui.Select{
		Label: "Select shell",
		Items: []string{
			"sh",
			"bash",
			"csh",
			"ksh",
			"zsh",
		},
	}

	return prompt.Run()
}
