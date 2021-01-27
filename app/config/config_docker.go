package config

// Docker project path

// GetProjectDockerPath gets project main container
func (c *Config) GetProjectDockerPath() string {
	return c.projectConfig.GetDockerProjectPath()
}

// SaveDockerProjectPath saves path to project in container into project file
func (c *Config) SaveDockerProjectPath(path string) (err error) {
	c.projectConfig.DockerProjectPath = path
	return c.saveProjectFile()
}

// Docker instance command

// SetDockerCommand define docker command
func (c *Config) SetDockerCommand(command string) error {
	c.globalConfig.SetDockerCommand(command)
	return c.fileSystem.SaveConfigFile(c.globalConfig, c.GetUserFile())
}

// GetDockerCommand gets the docker command
func (c *Config) GetDockerCommand() string {
	return c.globalConfig.GetDockerCommand()
}
