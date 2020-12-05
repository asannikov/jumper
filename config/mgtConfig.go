package config

import (
	"fmt"
)

// MgtConfig contains configuration from json config file
type MgtConfig struct {
	FileName        string
	FileUser        string
	globalSettings  GlobalConfig
	projectSettings ProjectSettings
	FileSystem      FileSystem
}

// GenerateConfig gets MgtConfig configuration
func (c *MgtConfig) GenerateConfig(properties ...func() error) error {
	if len(properties) == 0 {
		return fmt.Errorf("Define property evaluation")
	}

	var err error
	for _, p := range properties {
		err = p()

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *MgtConfig) handleConfig() error {
	pconfiguration := ProjectConfig{}
	if t, err := c.FileSystem.FileExists(c.FileName); err == nil && t == true {
		err := c.FileSystem.ReadConfigFile(c.FileName, &pconfiguration)

		if err != nil {
			return err
		}

		if pconfiguration.PhpContainer == "" {
			return fmt.Errorf("Php container missing in config file")
		}

		c.projectSettings.PhpContainer = pconfiguration.PhpContainer
	} else if err != nil {
		return err
	}

	gconfiguration := GlobalConfig{}
	if t, err := c.FileSystem.FileExists(c.FileUser); err == nil && t == true {
		err := c.FileSystem.ReadConfigFile(c.FileUser, &gconfiguration)
		if err != nil {
			return err
		}

		c.globalSettings = gconfiguration
	} else if err != nil {
		return err
	}

	return nil
}
