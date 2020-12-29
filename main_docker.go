package main

/**
 * @todo - improve:
 * dockerStartDialog type has privat func methods for unittesting.
 * It's not the best implementation, cause each funtion has no access to dockerStartDialog object pointer
 * and for that reason each such function has a parameter with the pointer within.
 */

import (
	"errors"
	"jumper/dialog"
	"jumper/docker"

	"github.com/docker/docker/client"
)

type dockerInstance interface {
	Run(string) error
	GetContanerList() ([]string, error)
	Stat() (string, error)
	InitClient() error
}

type dockerStartDialog struct {
	dialog            *dialog.Dialog
	docker            *docker.Docker
	dockerService     string
	stat              func(*dockerStartDialog) (string, error)
	initClient        func(*dockerStartDialog) error
	run               func(*dockerStartDialog) error
	startDockerDialog func(*dockerStartDialog) (string, error)
	getClient         func(*dockerStartDialog) *client.Client
	containerList     func(*dockerStartDialog) ([]string, error)
}

func getDockerStartDialog() *dockerStartDialog {
	cl := &dockerStartDialog{}
	cl.stat = func(cl *dockerStartDialog) (string, error) {
		return cl.docker.Stat()
	}
	cl.initClient = func(cl *dockerStartDialog) error {
		return cl.docker.InitClient()
	}

	cl.run = func(cl *dockerStartDialog) error {
		return cl.docker.Run(cl.dockerService)
	}

	cl.startDockerDialog = func(cl *dockerStartDialog) (string, error) {
		return cl.dialog.StartDocker()
	}

	cl.getClient = func(cl *dockerStartDialog) *client.Client {
		return cl.docker.GetClient()
	}

	cl.containerList = func(cl *dockerStartDialog) ([]string, error) {
		return cl.docker.GetContanerList()
	}

	return cl
}

func (dsd *dockerStartDialog) GetContainerList() ([]string, error) {
	var err error
	var apiVersion string

	if apiVersion, err = dsd.stat(dsd); err != nil && apiVersion != "" {
		return []string{}, err
	}

	var choice string

	if apiVersion == "" {
		choice, err = dsd.startDockerDialog(dsd)
	}

	if err != nil {
		return []string{}, err
	}

	if choice == "y" || choice == "Y" {
		err = dsd.run(dsd)
	} else {
		err = dsd.initClient(dsd)
	}

	if err != nil {
		return []string{}, err
	}

	if dsd.getClient(dsd) == nil {
		return []string{}, errors.New("This command requires Docker to be run. Please, start it first")
	}

	return dsd.containerList(dsd)
}

func (dsd *dockerStartDialog) setDialog(d dialogCommand) {
	dsd.dialog = d.(*dialog.Dialog)
}

func (dsd *dockerStartDialog) setDockerService(c string) {
	dsd.dockerService = c
}

func (dsd *dockerStartDialog) setDocker(d dockerInstance) {
	dsd.docker = d.(*docker.Docker)
}
