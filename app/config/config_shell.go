package config

// SaveShellCommand saves linux shell command
func (c *Config) SaveShellCommand(cmd string) (err error) {
	c.projectConfig.Shell = cmd
	return c.saveProjectFile()
}

// GetShell gets shell command
func (c *Config) GetShell() string {
	shell := c.projectConfig.GetShell()
	if shell == "" {
		return "sh"
	}
	return shell
}
