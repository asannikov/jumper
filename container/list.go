package container

// GO111MODULE=on go get github.com/docker/docker@v19.03.13

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// GetContanerList gets active docker container list
func GetContanerList() []string {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	var containerList []string
	for _, container := range containers {
		for _, name := range container.Names {
			containerList = append(containerList, strings.TrimLeft(name, "/"))
		}
	}

	return containerList
}
