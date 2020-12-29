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
	client            *client.Client
	exec              func(string, ...string) *exec.Cmd
	initClient        func(*Docker) error
	clientping        func(cli *client.Client) (types.Ping, error)
	ping              func(*Docker) (types.Ping, error)
	run               func(*Docker, string) error
	newClientWithOpts func(...client.Opt) (*client.Client, error)
}

// GetDockerInstance gets docker instance
func GetDockerInstance() *Docker {
	docker := &Docker{}

	docker.newClientWithOpts = client.NewClientWithOpts
	docker.exec = exec.Command

	docker.initClient = func(d *Docker) (err error) {
		var cli *client.Client
		if cli, err = d.newClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation()); err == nil {
			docker.client = cli
		}

		return err
	}

	docker.clientping = func(cli *client.Client) (types.Ping, error) {
		ctx := context.Background()
		ping, err := cli.Ping(ctx)

		if err != nil {
			return ping, err
		}

		return ping, nil
	}

	docker.ping = func(d *Docker) (types.Ping, error) {

		cli, err := d.newClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

		if err != nil {
			return types.Ping{}, err
		}

		return d.clientping(cli)
	}

	docker.run = func(d *Docker, cmd string) (err error) {
		if err = openDocker(cmd, d.exec); err != nil {
			return err
		}

		fmt.Print("Docker is loading")
		for range time.Tick(1 * time.Second) {
			connection, err := docker.Ping()
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

	return docker
}

// GetClient gets docker client
func (d *Docker) GetClient() *client.Client {
	return d.client
}

// InitClient gets docker client
func (d *Docker) InitClient() error {
	return d.initClient(d)
}

func getCommand(command string) (string, []string) {
	cmdSlice := strings.Split(command, " ")

	command = strings.Trim(cmdSlice[0], " ")
	args := []string{}

	for _, v := range cmdSlice[1:] {
		if v != "" {
			args = append(args, strings.Trim(v, " "))
		}
	}

	return command, args
}

// macos open --hide -a Docker
func openDocker(command string, ecmd func(string, ...string) *exec.Cmd) error {
	if len(strings.Trim(command, " ")) == 0 {
		return errors.New("Docker instance is empty. Please, define it in global config by starting any command")
	}

	command, args := getCommand(command)

	cmd := ecmd(command, args...)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Run starts docker service
func (d *Docker) Run(cmd string) (err error) {
	return d.run(d, cmd)
}

// Ping checks the docker instance state
func (d *Docker) Ping() (types.Ping, error) {
	return d.ping(d)
}

// Stat gets docker instance API version
func (d *Docker) Stat() (string, error) {
	connection, err := d.Ping()
	return connection.APIVersion, err
}
