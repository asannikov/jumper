package docker

// GO111MODULE=on go get github.com/docker/docker@v19.03.13

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
)

// GetContanerList get container list
func (d *Docker) GetContanerList() ([]string, error) {
	var err error
	var containers []types.Container

	if containers, err = d.GetClient().ContainerList(context.Background(), types.ContainerListOptions{}); err != nil {
		return nil, err
	}

	var containerList []string
	for _, container := range containers {
		for _, name := range container.Names {
			containerList = append(containerList, strings.TrimLeft(name, "/"))
		}
	}

	/*if len(containerList) == 0 {
		return nil, errors.New("No running projects. Start project using docker-compose up and then run again `jumper start`")
	}*/
	return containerList, nil
}
