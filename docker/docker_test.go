package docker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func TestInitClientCase1(t *testing.T) {
	docker := GetDockerInstance()

	docker.newClientWithOpts = func(opts ...client.Opt) (*client.Client, error) {
		return &client.Client{}, nil
	}

	err := docker.InitClient()

	assert.Nil(t, err)
	assert.Equal(t, &client.Client{}, docker.GetClient())
}

func TestInitClientCase2(t *testing.T) {
	docker := GetDockerInstance()

	docker.newClientWithOpts = func(opts ...client.Opt) (*client.Client, error) {
		return nil, fmt.Errorf("unable to verify TLS configuration, invalid transport")
	}

	err := docker.InitClient()

	assert.EqualError(t, err, "unable to verify TLS configuration, invalid transport")
	assert.Nil(t, docker.GetClient())
}

func TestPingCase1(t *testing.T) {
	docker := GetDockerInstance()

	docker.newClientWithOpts = func(opts ...client.Opt) (*client.Client, error) {
		return nil, fmt.Errorf("unable to verify TLS configuration, invalid transport")
	}

	ping, err := docker.Ping()

	assert.EqualError(t, err, "unable to verify TLS configuration, invalid transport")
	assert.Equal(t, types.Ping{}, ping)
}
