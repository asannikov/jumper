package main

import (
	"errors"
	"jumper/dialog"
	"jumper/docker"
)

type dockerInstance interface {
	Run() error
	GetContanerList() ([]string, error)
	Stat() (string, error)
	InitClient() error
}

type dockerStartDialog struct {
	dialog *dialog.Dialog
	docker *docker.Docker
}

func (dsd *dockerStartDialog) GetContainerList() ([]string, error) {
	var err error
	var apiVersion string

	if apiVersion, err = dsd.docker.Stat(); err != nil && apiVersion != "" {
		return nil, err
	}

	var choice string

	if apiVersion == "" {
		choice, err = dsd.dialog.StartDocker()
	}

	if err != nil {
		return nil, err
	}

	if choice == "y" || choice == "Y" {
		err = dsd.docker.Run()
	} else {
		err = dsd.docker.InitClient()
	}

	if err != nil {
		return nil, err
	}

	if dsd.docker.GetClient() == nil {
		return nil, errors.New("This command requires Docker to be run. Please, start it first")
	}

	return dsd.docker.GetContanerList()
}

func (dsd *dockerStartDialog) setDialog(d dialogCommand) {
	dsd.dialog = d.(*dialog.Dialog)
}

func (dsd *dockerStartDialog) setDocker(d dockerInstance) {
	dsd.docker = d.(*docker.Docker)
}
