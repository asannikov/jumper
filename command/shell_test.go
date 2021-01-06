package command

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testShellConfig struct {
	shell string
	save  error
}

func (s *testShellConfig) GetShell() string {
	return s.shell
}

func (s *testShellConfig) SaveShellCommand(v string) error {
	return s.save
}

type testShellDialog struct {
	shell string
	err   error
}

func (s *testShellDialog) DockerShell() (int, string, error) {
	return 0, s.shell, s.err
}

func TestDefineShellTypeCase1(t *testing.T) {
	cfg := &testShellConfig{}
	d := &testShellDialog{
		shell: "bash",
		err:   nil,
	}

	assert.Nil(t, defineShellType(cfg, d))
}

func TestDefineShellTypeCase2(t *testing.T) {
	cfg := &testShellConfig{}
	d := &testShellDialog{
		shell: "",
		err:   errors.New("Something goes wrong. Shell was not set"),
	}

	assert.EqualError(t, defineShellType(cfg, d), "Something goes wrong. Shell was not set")
}

func TestDefineShellTypeCase3(t *testing.T) {
	cfg := &testShellConfig{
		save: errors.New("cannot save shell"),
	}
	d := &testShellDialog{
		shell: "bash",
		err:   nil,
	}

	assert.EqualError(t, defineShellType(cfg, d), "cannot save shell")
}
