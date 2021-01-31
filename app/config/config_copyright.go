package config

// Copyright

// EnableCopyright Enable copyright output
func (c *Config) EnableCopyright() error {
	c.globalConfig.EnableCopyright()
	return c.fileSystem.SaveConfigFile(c.globalConfig, c.GetUserFile())
}

// DisableCopyright Disable copyright output
func (c *Config) DisableCopyright() error {
	c.globalConfig.DisableCopyright()
	return c.fileSystem.SaveConfigFile(c.globalConfig, c.GetUserFile())
}

// ShowCopyrightText check the status of copyright output
func (c *Config) ShowCopyrightText() bool {
	return c.globalConfig.ShowCopyrightText()
}
