package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tkanos/gonfig"
)

// MgtConfig contains configuration from json config file
type MgtConfig struct {
	FileName string
}

func (c *MgtConfig) saveConfigFile(data Configuration) error {
	file, _ := json.MarshalIndent(data, "", " ")
	err := ioutil.WriteFile(c.FileName, file, 0644)
	return err
}

// v checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func (c *MgtConfig) fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetConfig gets MgtConfig configuration
func (c *MgtConfig) GetConfig() (Configuration, error) {

	configuration := Configuration{}

	if fileExists(c.FileName) {
		err := gonfig.GetConf(c.FileName, &configuration)
		return configuration, err
	}

	DLG := dialogInternal{
		setPhpContaner: selectPhpContainer,
	}

	_, containerPhp, err := DLG.setPhpContaner()

	if err != nil {
		return configuration, err
	}

	err = saveConfigFile(Configuration{
		PhpContainer: containerPhp,
	})

	if err != nil {
		return configuration, err
	}

	if fileExists(c.FileName) {
		return getConfigFile()
	}

	return configuration, fmt.Errorf("Config file %s does not exist ", c.FileName)
}

// GetPhpContainer gets PHP container name
func (c *MgtConfig) GetPhpContainer() (string, error) {
	cfg, err := c.GetConfig()

	if err == nil {
		return cfg.PhpContainer, nil
	}

	return "", err
}
