package command

import (
	"errors"
)

type projectConfig interface {
	GetProjectMainContainer() string
	GetStartCommand() string
	GetProjectDockerPath() string
	GetXDebugCliIniPath() string
	GetXDebugFpmIniPath() string
	GetXDebugConifgLocaton() string
	SaveContainerNameToProjectConfig(string) error
	SaveStartCommandToProjectConfig(string) error
	SaveDockerProjectPath(string) error
	SaveDockerCliXdebugIniFilePath(string) error
	SaveDockerFpmXdebugIniFilePath(string) error
	SaveXDebugConifgLocaton(string) error
}

type dialog interface {
	SetMainContaner([]string) (int, string, error)
	StartCommand() (string, error)
	StartDocker() (string, error)
	DockerService() (string, error)
	DockerProjectPath(string) (string, error)
	DockerCliXdebugIniFilePath(string) (string, error)
	DockerFpmXdebugIniFilePath(string) (string, error)
	XDebugConfigLocation() (int, string, error)
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

func defineProjectDockerPath(cfg projectConfig, d dialog, defaultPath string) (err error) {
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

func defineCliXdebugIniFilePath(cfg projectConfig, d dialog, defaultPath string) (err error) {
	if cfg.GetXDebugCliIniPath() == "" {
		var path string
		if path, err = d.DockerCliXdebugIniFilePath(defaultPath); err != nil {
			return err
		}

		if path == "" {
			return errors.New("Cli Xdebug ini file path is empty")
		}

		return cfg.SaveDockerCliXdebugIniFilePath(path)
	}

	return nil
}

func defineFpmXdebugIniFilePath(cfg projectConfig, d dialog, defaultPath string) (err error) {
	if cfg.GetXDebugFpmIniPath() == "" {
		var path string
		if path, err = d.DockerFpmXdebugIniFilePath(defaultPath); err != nil {
			return err
		}

		if path == "" {
			return errors.New("Fpm Xdebug ini file path is empty")
		}

		return cfg.SaveDockerFpmXdebugIniFilePath(path)
	}

	return nil
}

func defineXdebugIniFileLocation(cfg projectConfig, d dialog) (err error) {
	if cfg.GetXDebugConifgLocaton() == "" {
		var path string

		if _, path, err = d.XDebugConfigLocation(); err != nil {
			return err
		}

		if path == "" {
			return errors.New("Xdebug config file locaton cannot be empty")
		}

		return cfg.SaveXDebugConifgLocaton(path)
	}

	return nil
}
