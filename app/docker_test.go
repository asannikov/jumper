package app

import (
	"errors"
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

type testDockerInstance struct {
	stat             func() (string, error)
	run              func(string) error
	initClient       func() error
	getContainerList func() ([]string, error)
}

func (d *testDockerInstance) Run(cmd string) error {
	return d.run(cmd)
}

func (d *testDockerInstance) GetContanerList() ([]string, error) {
	return d.getContainerList()
}

func (d *testDockerInstance) Stat() (string, error) {
	return d.stat()
}

func (d *testDockerInstance) InitClient() error {
	return d.initClient()
}

func TestGetContainerListCase1(t *testing.T) {
	cl := &dockerStartDialog{}
	cl.stat = func(cl *dockerStartDialog) (string, error) {
		return "1.0", errors.New("Ping error")
	}

	cl.setDockerService("run docker")

	list, err := cl.GetContainerList()

	assert.Equal(t, []string{}, list)
	assert.EqualError(t, err, "Ping error")
}

func TestGetContainerListCase2(t *testing.T) {
	cl := &dockerStartDialog{}
	cl.stat = func(cl *dockerStartDialog) (string, error) {
		return "", nil
	}
	cl.startDockerDialog = func(cl *dockerStartDialog) (string, error) {
		return "", errors.New("dialog error")
	}

	cl.setDockerService("run docker")

	list, err := cl.GetContainerList()

	assert.Equal(t, []string{}, list)
	assert.EqualError(t, err, "dialog error")
}

func TestGetContainerListCase3(t *testing.T) {
	cl := &dockerStartDialog{}
	cl.stat = func(cl *dockerStartDialog) (string, error) {
		return "", nil
	}
	cl.startDockerDialog = func(cl *dockerStartDialog) (string, error) {
		return "y", nil
	}
	cl.run = func(cl *dockerStartDialog) error {
		return errors.New("command error")
	}

	cl.setDockerService("run docker")

	list, err := cl.GetContainerList()

	assert.Equal(t, []string{}, list)
	assert.EqualError(t, err, "command error")
}

func TestGetContainerListCase4(t *testing.T) {
	cl := &dockerStartDialog{}
	cl.stat = func(cl *dockerStartDialog) (string, error) {
		return "", nil
	}
	cl.startDockerDialog = func(cl *dockerStartDialog) (string, error) {
		return "Y", nil
	}
	cl.run = func(cl *dockerStartDialog) error {
		return errors.New("command error")
	}

	cl.setDockerService("run docker")

	list, err := cl.GetContainerList()

	assert.Equal(t, []string{}, list)
	assert.EqualError(t, err, "command error")
}

func TestGetContainerListCase5(t *testing.T) {
	cl := &dockerStartDialog{}
	cl.stat = func(cl *dockerStartDialog) (string, error) {
		return "", nil
	}
	cl.startDockerDialog = func(cl *dockerStartDialog) (string, error) {
		return "N", nil
	}
	cl.run = func(cl *dockerStartDialog) error {
		return errors.New("command error")
	}
	cl.initClient = func(cl *dockerStartDialog) error {
		return errors.New("init client error")
	}
	cl.setDockerService("run docker")

	list, err := cl.GetContainerList()

	assert.Equal(t, []string{}, list)
	assert.EqualError(t, err, "init client error")
}

func TestGetContainerListCase6(t *testing.T) {
	cl := &dockerStartDialog{}
	cl.stat = func(cl *dockerStartDialog) (string, error) {
		return "", nil
	}
	cl.startDockerDialog = func(cl *dockerStartDialog) (string, error) {
		return "N", nil
	}
	cl.run = func(cl *dockerStartDialog) error {
		return errors.New("command error")
	}
	cl.initClient = func(cl *dockerStartDialog) error {
		return nil
	}
	cl.getClient = func(cl *dockerStartDialog) *client.Client {
		return nil
	}
	cl.setDockerService("run docker")

	list, err := cl.GetContainerList()

	assert.Equal(t, []string{}, list)
	assert.EqualError(t, err, "This command requires Docker to be run. Please, start it first")
}

func TestGetContainerListCase7(t *testing.T) {
	cl := &dockerStartDialog{}
	cl.stat = func(cl *dockerStartDialog) (string, error) {
		return "", nil
	}
	cl.startDockerDialog = func(cl *dockerStartDialog) (string, error) {
		return "N", nil
	}
	cl.run = func(cl *dockerStartDialog) error {
		return errors.New("command error")
	}
	cl.initClient = func(cl *dockerStartDialog) error {
		return nil
	}
	cl.getClient = func(cl *dockerStartDialog) *client.Client {
		return &client.Client{}
	}
	cl.containerList = func(cl *dockerStartDialog) ([]string, error) {
		return []string{"test_container"}, nil
	}
	cl.setDockerService("run docker")

	list, err := cl.GetContainerList()

	assert.Equal(t, []string{"test_container"}, list)
	assert.Nil(t, err)
}
