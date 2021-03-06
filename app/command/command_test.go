package command

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testConfigCommand struct {
	getMainContainerUser             string
	getProjectMainContainer          string
	getProjectDockerPath             string
	saveContainerUserToProjectConfig error
	saveContainerNameToProjectConfig error
	saveDockerProjectPath            error
}

func (c *testConfigCommand) SaveContainerNameToProjectConfig(string) error {
	return c.saveContainerNameToProjectConfig
}
func (c *testConfigCommand) SaveContainerUserToProjectConfig(string) error {
	return c.saveContainerUserToProjectConfig
}
func (c *testConfigCommand) SaveDockerProjectPath(s string) error {
	return c.saveDockerProjectPath
}
func (c *testConfigCommand) GetMainContainerUser() string {
	return c.getMainContainerUser
}
func (c *testConfigCommand) GetProjectMainContainer() string {
	return c.getProjectMainContainer
}
func (c *testConfigCommand) GetProjectDockerPath() string {
	return c.getProjectDockerPath
}

type testDialogCommand struct {
	setMainContainer    func([]string) (int, string, error)
	dockerProjectPath   func(string) (string, error)
	setMainContanerUser func() (string, error)
}

func (d *testDialogCommand) SetMainContaner(l []string) (int, string, error) {
	return d.setMainContainer(l)
}
func (d *testDialogCommand) SetMainContanerUser() (string, error) {
	return d.setMainContanerUser()
}

func (d *testDialogCommand) DockerProjectPath(s string) (string, error) {
	return d.dockerProjectPath(s)
}

func TestDefineProjectMainContainerCase1(t *testing.T) {
	cfg := &testConfigCommand{
		getProjectMainContainer: "main_container",
	}
	dlg := &testDialogCommand{}

	assert.Nil(t, defineProjectMainContainer(cfg, dlg, []string{}))
}

func TestDefineProjectMainContainerCase2(t *testing.T) {
	cfg := &testConfigCommand{
		getProjectMainContainer: "",
	}
	dlg := &testDialogCommand{
		setMainContainer: func(l []string) (int, string, error) {
			return 0, "", errors.New("dialog container error")
		},
	}

	assert.EqualError(t, defineProjectMainContainer(cfg, dlg, []string{"container"}), "dialog container error")
}

func TestDefineProjectMainContainerCase3(t *testing.T) {
	cfg := &testConfigCommand{
		getProjectMainContainer: "",
	}
	dlg := &testDialogCommand{
		setMainContainer: func(l []string) (int, string, error) {
			return 0, "", nil
		},
	}

	assert.EqualError(t, defineProjectMainContainer(cfg, dlg, []string{"container"}), "Container name is empty. Set the container name")
}

func TestDefineProjectMainContainerCase4(t *testing.T) {
	cfg := &testConfigCommand{
		getProjectMainContainer:          "",
		saveContainerNameToProjectConfig: errors.New("saveContainerNameToProjectConfig error"),
	}
	dlg := &testDialogCommand{
		setMainContainer: func(l []string) (int, string, error) {
			return 0, "main_container", nil
		},
	}

	assert.EqualError(t, defineProjectMainContainer(cfg, dlg, []string{"container"}), "saveContainerNameToProjectConfig error")
}

func TestDefineProjectMainContainerCase5(t *testing.T) {
	cfg := &testConfigCommand{
		getProjectMainContainer:          "",
		saveContainerNameToProjectConfig: nil,
	}
	dlg := &testDialogCommand{
		setMainContainer: func(l []string) (int, string, error) {
			return 0, "main_container", nil
		},
	}

	assert.Nil(t, defineProjectMainContainer(cfg, dlg, []string{"container"}))
}

func TestDefineProjectMainContainerUserCase1(t *testing.T) {
	cfg := &testConfigCommand{
		getMainContainerUser: "user",
	}
	dlg := &testDialogCommand{}

	assert.Nil(t, defineProjectMainContainerUser(cfg, dlg))
}

func TestDefineProjectMainContainerUserCase2(t *testing.T) {
	cfg := &testConfigCommand{
		getMainContainerUser: "",
	}
	dlg := &testDialogCommand{
		setMainContanerUser: func() (string, error) {
			return "", errors.New("dialog user error")
		},
	}

	assert.EqualError(t, defineProjectMainContainerUser(cfg, dlg), "dialog user error")
}

func TestDefineProjectMainContainerUserCase3(t *testing.T) {
	cfg := &testConfigCommand{
		getMainContainerUser: "",
	}
	dlg := &testDialogCommand{
		setMainContanerUser: func() (string, error) {
			return "", nil
		},
	}

	assert.EqualError(t, defineProjectMainContainerUser(cfg, dlg), "Container user name is empty. Set the user name")
}

func TestDefineProjectMainContainerUserCase4(t *testing.T) {
	cfg := &testConfigCommand{
		getMainContainerUser:             "",
		saveContainerUserToProjectConfig: errors.New("saveContainerUserToProjectConfig error"),
	}
	dlg := &testDialogCommand{
		setMainContanerUser: func() (string, error) {
			return "main_container", nil
		},
	}

	assert.EqualError(t, defineProjectMainContainerUser(cfg, dlg), "saveContainerUserToProjectConfig error")
}

func TestDefineProjectMainContainerUserCase5(t *testing.T) {
	cfg := &testConfigCommand{
		getProjectMainContainer:          "",
		saveContainerUserToProjectConfig: nil,
	}
	dlg := &testDialogCommand{
		setMainContanerUser: func() (string, error) {
			return "main_container", nil
		},
	}

	assert.Nil(t, defineProjectMainContainerUser(cfg, dlg))
}

func TestDefineProjectDockerPathCase1(t *testing.T) {
	cfg := &testConfigCommand{
		getProjectDockerPath: "docker_path",
	}
	dlg := &testDialogCommand{}

	assert.Nil(t, defineProjectDockerPath(cfg, dlg, "default_path"))
}

func TestDefineProjectDockerPathCase2(t *testing.T) {
	cfg := &testConfigCommand{
		getProjectDockerPath: "",
	}
	dlg := &testDialogCommand{
		dockerProjectPath: func(s string) (string, error) {
			return "", errors.New("dialog docker path error")
		},
	}

	assert.EqualError(t, defineProjectDockerPath(cfg, dlg, "default_path"), "dialog docker path error")
}

func TestDefineProjectDockerPathCase3(t *testing.T) {
	cfg := &testConfigCommand{
		getProjectDockerPath: "",
	}
	dlg := &testDialogCommand{
		dockerProjectPath: func(s string) (string, error) {
			return "", nil
		},
	}

	assert.EqualError(t, defineProjectDockerPath(cfg, dlg, "default_path"), "Container path is empty. Set the path to project in container")
}

func TestDefineProjectDockerPathCase4(t *testing.T) {
	cfg := &testConfigCommand{
		getProjectDockerPath:  "",
		saveDockerProjectPath: errors.New("saveDockerProjectPath error"),
	}
	dlg := &testDialogCommand{
		dockerProjectPath: func(s string) (string, error) {
			return "docker/path", nil
		},
	}

	assert.EqualError(t, defineProjectDockerPath(cfg, dlg, "default_path"), "saveDockerProjectPath error")
}

func TestDefineProjectDockerPathCase5(t *testing.T) {
	cfg := &testConfigCommand{
		getProjectDockerPath:  "",
		saveDockerProjectPath: nil,
	}
	dlg := &testDialogCommand{
		dockerProjectPath: func(s string) (string, error) {
			return "docker/path", nil
		},
	}

	assert.Nil(t, defineProjectDockerPath(cfg, dlg, "default_path"))
}
