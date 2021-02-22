package config

// main container

// SaveContainerUserToProjectConfig saves container name into project file
func (c *Config) SaveContainerUserToProjectConfig(cu string) (err error) {
	c.projectConfig.ContainerUser = cu
	return c.saveProjectFile()
}

// GetMainContainerUser gets project main container
func (c *Config) GetMainContainerUser() string {
	return c.projectConfig.GetMainContainerUser()
}
