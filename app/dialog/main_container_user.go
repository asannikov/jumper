package dialog

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// SetMainContanerUser sets main container name
func (d *Dialog) SetMainContanerUser() (string, error) {
	return d.setMainContanerUser()
}

func setMainContanerUser() (string, error) {
	validate := func(u string) error {
		if u == "" {
			return fmt.Errorf("User name cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Add main docker container user",
		Validate: validate,
		Default:  "root",
	}

	return prompt.Run()
}
