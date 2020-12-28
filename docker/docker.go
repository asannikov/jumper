package docker

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Docker is a wrapper for docker SDK
type Docker struct {
	client *client.Client
}

// GetClient gets docker client
func (d *Docker) GetClient() *client.Client {
	return d.client
}

// InitClient gets docker client
func (d *Docker) InitClient() (err error) {
	if cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation()); err == nil {
		d.client = cli
	}
	return err
}

// macos open --hide -a Docker
func openDocker(command string) error {

	cmdSlice := strings.Split(command, " ")

	if len(cmdSlice) == 0 {
		return errors.New("Docker instance is empty. Please, define it")
	}

	command = strings.Trim(cmdSlice[0], " ")
	args := []string{}

	for _, v := range cmdSlice[1:] {
		args = append(args, strings.Trim(v, " "))
	}

	cmd := exec.Command(command, args...)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Run starts docker service
func (d *Docker) Run(cmd string) (err error) {
	if err = openDocker(cmd); err != nil {
		return err
	}

	fmt.Print("Docker is loading")
	for range time.Tick(1 * time.Second) {
		connection, err := d.Ping()
		fmt.Print("...")

		if connection.APIVersion == "" {
			continue
		}

		if err != nil {
			return err
		}

		fmt.Println("\nDocker is running")
		break
	}

	return d.InitClient()
}

// Ping checks the docker instance state
func (d *Docker) Ping() (types.Ping, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return types.Ping{}, err
	}

	ping, err := cli.Ping(ctx)

	if err != nil {
		return ping, err
	}

	return ping, nil
}

// Stat gets docker instance API version
func (d *Docker) Stat() (string, error) {
	connection, err := d.Ping()
	return connection.APIVersion, err
}
