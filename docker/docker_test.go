package docker

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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

func TestPingCase2(t *testing.T) {
	docker := GetDockerInstance()

	docker.newClientWithOpts = func(opts ...client.Opt) (*client.Client, error) {
		return &client.Client{}, nil
	}

	docker.clientping = func(cli *client.Client) (types.Ping, error) {
		return types.Ping{}, fmt.Errorf("client ping error")
	}

	ping, err := docker.Ping()

	assert.EqualError(t, err, "client ping error")
	assert.Equal(t, types.Ping{}, ping)
}

func TestPingCase3(t *testing.T) {
	docker := GetDockerInstance()

	docker.newClientWithOpts = func(opts ...client.Opt) (*client.Client, error) {
		return &client.Client{}, nil
	}

	docker.clientping = func(cli *client.Client) (types.Ping, error) {
		return types.Ping{}, nil
	}

	ping, err := docker.Ping()

	assert.Nil(t, err)
	assert.Equal(t, types.Ping{}, ping)
}

func TestGetCommandCase1(t *testing.T) {
	command, args := getCommand("macos open --hide -a Docker")
	assert.Equal(t, "macos", command)
	assert.Equal(t, []string{"open", "--hide", "-a", "Docker"}, args)
}

func TestGetCommandCase2(t *testing.T) {
	command, args := getCommand("macos  open --hide -a  docker")
	assert.Equal(t, "macos", command)
	assert.Equal(t, []string{"open", "--hide", "-a", "docker"}, args)
}

func TestGetCommandCase3(t *testing.T) {
	command, args := getCommand("macos  open - -hide -a  docker")
	assert.Equal(t, "macos", command)
	assert.Equal(t, []string{"open", "-", "-hide", "-a", "docker"}, args)
}

func TestGetCommandCase4(t *testing.T) {
	command, args := getCommand("")
	assert.Equal(t, "", command)
	assert.Equal(t, []string{}, args)
}

func TestGetCommandCase5(t *testing.T) {
	command, args := getCommand(" ")
	assert.Equal(t, "", command)
	assert.Equal(t, []string{}, args)
}

const stub = "test_file\ntest_file2\n"

var testCase string

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	tc := "TEST_CASE=" + testCase
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1", tc}

	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}
	switch os.Getenv("TEST_CASE") {
	case "case1":
		fmt.Fprint(os.Stdout, stub)
	}
}

func TestOpenDockerCase1(t *testing.T) {
	testCase = "case1"
	execCommand := fakeExecCommand
	defer func() { execCommand = exec.Command }()
	err := openDocker("open -a docker", execCommand)
	assert.Nil(t, err)
}

func TestRunCase1(t *testing.T) {
	docker := GetDockerInstance()

	docker.exec = fakeExecCommand

	docker.initClient = func() error {
		return nil
	}

	docker.ping = func() (types.Ping, error) {
		t := types.Ping{
			APIVersion: "1.0",
		}
		return t, nil
	}

	err := docker.Run("open -a docker")

	assert.Nil(t, err)
}

func TestRunCase2(t *testing.T) {
	docker := GetDockerInstance()

	docker.exec = fakeExecCommand

	docker.initClient = func() error {
		return nil
	}

	docker.ping = func() (types.Ping, error) {
		t := types.Ping{
			APIVersion: "1.0",
		}
		return t, errors.New("client ping error")
	}

	err := docker.Run("open -a docker")

	assert.EqualError(t, err, "client ping error")
}
