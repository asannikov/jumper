package config

// main start command

// GetStartCommand gets start command
func (c *Config) GetStartCommand() string {
	return c.projectConfig.GetStartCommand()
}

// SaveStartCommandToProjectConfig saves container name into project file
func (c *Config) SaveStartCommandToProjectConfig(cmd string) (err error) {
	c.projectConfig.StartCommand = cmd
	return c.saveProjectFile()
}
