package command

import (
	"errors"
)

type containerlist interface {
	GetContainerList() ([]string, error)
}

type projectMainContainerProjectConfig interface {
	SaveContainerNameToProjectConfig(string) error
	GetProjectMainContainer() string
}

type defineProjectMainContainerDialog interface {
	SetMainContaner([]string) (int, string, error)
}

func defineProjectMainContainer(cfg projectMainContainerProjectConfig, d defineProjectMainContainerDialog, containerlist []string) (err error) {
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

type dockerPathProjectConfig interface {
	SaveDockerProjectPath(string) error
	GetProjectDockerPath() string
}

type defineProjectDockerPathDialog interface {
	DockerProjectPath(string) (string, error)
}

func defineProjectDockerPath(cfg dockerPathProjectConfig, d defineProjectDockerPathDialog, defaultPath string) (err error) {
	if cfg.GetProjectDockerPath() == "" {
		var path string
		if path, err = d.DockerProjectPath(defaultPath); err != nil {
			return err
		}

		if path == "" {
			return errors.New("Container path is empty. Set the path to project in container")
		}

		return cfg.SaveDockerProjectPath(path)
	}

	return nil
}
