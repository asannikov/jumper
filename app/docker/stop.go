package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
)

// StopContainers stop containers by filter
func (d *Docker) StopContainers() func([]string) error {
	return func(filter []string) (err error) {
		var containers []types.Container

		ctx := context.Background()

		if containers, err = d.GetClient().ContainerList(ctx, types.ContainerListOptions{}); err != nil {
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

			if err := d.GetClient().ContainerStop(ctx, container.ID, nil); err != nil {
				return err
			}

			fmt.Println("Success")
		}

		return nil
	}
}
