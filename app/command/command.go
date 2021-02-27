package command

import (
	"errors"
)

// ExecOptions defines exec options
type ExecOptions struct {
	command    string
	args       []string
	user       string
	tty        bool
	detach     bool
	workingDir string
}

// GetCommand gets main command, ie docker
func (eo *ExecOptions) GetCommand() string {
	return eo.command
}

// GetArgs gets arguments for the command
func (eo *ExecOptions) GetArgs() []string {
	return eo.args
}

// GetUser gets user to use in docker container
func (eo *ExecOptions) GetUser() string {
	return eo.user
}

// GetTty returns tty mode status
func (eo *ExecOptions) GetTty() bool {
	return eo.tty
}

// GetDetach returns detach mode status
func (eo *ExecOptions) GetDetach() bool {
	return eo.detach
}

// GetWorkingDir returns working directory
func (eo *ExecOptions) GetWorkingDir() string {
	return eo.workingDir
}

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

type projectMainContainerProjectUserConfig interface {
	SaveContainerUserToProjectConfig(string) error
	GetMainContainerUser() string
}

type defineProjectMainContainerUserDialog interface {
	SetMainContanerUser() (string, error)
}

func defineProjectMainContainerUser(cfg projectMainContainerProjectUserConfig, d defineProjectMainContainerUserDialog) (err error) {
	if cfg.GetMainContainerUser() == "" {
		user, err := d.SetMainContanerUser()

		if err != nil {
			return err
		}

		if user == "" {
			return errors.New("Container user name is empty. Set the user name")
		}

		return cfg.SaveContainerUserToProjectConfig(user)
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
