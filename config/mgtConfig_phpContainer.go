package config

// EvaluatePhpContainer helps to find/define php container
func (c *MgtConfig) EvaluatePhpContainer(spc func() (int, string, error)) error {
	configuration := ProjectConfig{}

	if t, err := c.FileSystem.FileExists(c.FileName); err == nil && t == true {
		err := c.FileSystem.ReadConfigFile(c.FileName, &configuration)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if configuration.PhpContainer == "" {
		_, containerPhp, err := spc()

		if err != nil {
			return err
		}

		configuration.PhpContainer = containerPhp
		err = c.FileSystem.SaveConfigFile(configuration, c.FileName)

		if err != nil {
			return err
		}
	}

	return nil
}

// GetPhpContainer gets PHP cosntainer name
func (c *MgtConfig) GetPhpContainer(gf func() error) (string, error) {
	err := gf()
	if err != nil {
		return "", err
	}

	err = c.handleConfig()

	if err != nil {
		return "", err
	}

	return c.projectSettings.PhpContainer, nil
}
