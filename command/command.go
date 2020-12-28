package command

import (
	"errors"
)

type projectConfig interface {
	GetProjectMainContainer() string
	GetStartCommand() string
	SaveContainerNameToProjectConfig(string) error
	SaveStartCommandToProjectConfig(string) error
}

type dialog interface {
	SetMainContaner([]string) (int, string, error)
	SetStartCommand() (string, error)
	StartDocker() (string, error)
}

type containerlist interface {
	GetContainerList() ([]string, error)
}

func defineProjectMainContainer(cfg projectConfig, d dialog, containerlist []string) (err error) {
	if cfg.GetProjectMainContainer() == "" {
		_, container, err := d.SetMainContaner(containerlist)

		if err != nil {
			return err
		}

		if container == "" {
			return errors.New("Container name is empty. Set the container name")
		}

		return cfg.SaveContainerNameToProjectConfig(container)
	}

	return nil
}
