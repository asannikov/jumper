package config

/*
// EvaluatePhpContainer helps to find/define php container
func (c *MgtConfig) EvaluatePhpContainer(spc func() (int, string, error)) error {
	configuration := ProjectConfig{}

	err := c.fileSystem.ReadConfigFile(c.FileName, &configuration)
	if err != nil {
		return err
	}

	if configuration.PhpContainer == "" {
		_, containerPhp, err := spc()

		if err != nil {
			return err
		}

		configuration.PhpContainer = containerPhp
		err = c.fileSystem.SaveConfigFile(configuration, c.FileName)

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

	err = c.loadProjectConfigFromJSON()

	if err != nil {
		return "", err
	}

	if c.projectSettings.PhpContainer == "" {
		return "", fmt.Errorf("Php container missing in config file")
	}

	return c.projectSettings.PhpContainer, nil
}
*/
