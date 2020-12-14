Add new project fields here `config/ProjectConfig.go` 
and add the related method to `config/config.go`, ie:
```
// SaveContainerNameToProjectConfig saves container name into project file
func (c *Config) SaveContainerNameToProjectConfig(cn string) (err error) {
	c.projectConfig.MainContainer = cn
	return c.fileSystem.SaveConfigFile(c.projectConfig, c.getProjectFile())
}

// GetProjectMainContainer gets project main container
func (c *Config) GetProjectMainContainer() string {
	return c.projectConfig.GetMainContainer()
}
```