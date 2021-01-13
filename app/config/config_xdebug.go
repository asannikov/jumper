package config

// Xdebug

// GetXDebugConfigLocaton gets cli xdebug ini file path
func (c *Config) GetXDebugConfigLocaton() string {
	return c.projectConfig.GetXDebugConfigLocaton()
}

// GetXDebugCliIniPath gets cli xdebug ini file path
func (c *Config) GetXDebugCliIniPath() string {
	return c.projectConfig.GetXDebugCliIniPath()
}

// GetXDebugFpmIniPath gets fpm xdebug ini file path
func (c *Config) GetXDebugFpmIniPath() string {
	return c.projectConfig.GetXDebugFpmIniPath()
}

// SaveXDebugConifgLocaton saves xdebug file location
func (c *Config) SaveXDebugConifgLocaton(path string) (err error) {
	c.projectConfig.XDebugLocation = path
	return c.saveProjectFile()
}

// SaveDockerCliXdebugIniFilePath saves xdebug cli ini file path into project file
func (c *Config) SaveDockerCliXdebugIniFilePath(path string) (err error) {
	c.projectConfig.XDebugCliIniPath = path
	return c.saveProjectFile()
}

// SaveDockerFpmXdebugIniFilePath saves xdebug fpm ini file path into project file
func (c *Config) SaveDockerFpmXdebugIniFilePath(path string) (err error) {
	c.projectConfig.XDebugFpmIniPath = path
	return c.saveProjectFile()
}
