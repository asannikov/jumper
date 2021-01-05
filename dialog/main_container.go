package dialog

import "github.com/manifoldco/promptui"

// SetMainContaner sets main container name
func (d *Dialog) SetMainContaner(cl []string) (int, string, error) {
	return d.setMainContaner(cl)
}

func setMainContaner(containers []string) (int, string, error) {
	prompt := promptui.Select{
		Label: "Select main container",
		Items: containers,
	}

	return prompt.Run()
}
