package dialog

import "github.com/manifoldco/promptui"

// SelectProject call project dropdown
func (d *Dialog) SelectProject(list []string) (int, string, error) {
	return d.setSelectProject(list)
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
