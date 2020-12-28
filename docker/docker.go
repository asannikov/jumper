package docker

import (
	"context"
	"fmt"
	"os/exec"
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

// Run starts docker service
func (d *Docker) Run() (err error) {
	cmd := exec.Command("open", "--hide", "-a", "Docker")
	if err = cmd.Run(); err != nil {
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
