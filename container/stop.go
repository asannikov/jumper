package container

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// StopContainers stops docker containers
func StopContainers() func([]string) error {
	return func(filter []string) (err error) {
		ctx := context.Background()

		var cli *client.Client
		var containers []types.Container

		if cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation()); err != nil {
			return err
		}

		if containers, err = cli.ContainerList(ctx, types.ContainerListOptions{}); err != nil {
			return err
		}

		for _, container := range containers {

			if len(filter) > 0 {
				found := false
				for _, name := range container.Names {
					for _, c := range filter {
						if c == strings.TrimLeft(string(name), "/") {
							found = true
							break
						}
					}
				}

				if !found {
					continue
				}
			}

			fmt.Print("\nStopping container ", container.ID[:10], "... ")

			for _, name := range container.Names {
				fmt.Print(name, " ... ")
			}

			if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
				return err
			}

			fmt.Println("Success")
		}

		return nil
	}
}
