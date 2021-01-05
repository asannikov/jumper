package command

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testComposerConfig struct {
	mainContainer string
}

func (tc *testComposerConfig) GetProjectMainContainer() string {
	return tc.mainContainer
}

func (tc *testComposerConfig) GetStartCommand() string {
	return ""
}

func (tc *testComposerConfig) GetProjectDockerPath() string {
	return ""
}

func (tc *testComposerConfig) SaveContainerNameToProjectConfig(container string) error {
	return nil
}

func (tc *testComposerConfig) SaveStartCommandToProjectConfig(c string) error {
	return nil
}

func (tc *testComposerConfig) SaveDockerProjectPath(c string) error {
	return nil
}

func (tc *testComposerConfig) SaveDockerCliXdebugIniFilePath(c string) error {
	return nil
}

func (tc *testComposerConfig) SaveDockerFpmXdebugIniFilePath(c string) error {
	return nil
}

func (tc *testComposerConfig) SaveXDebugConifgLocaton(c string) error {
	return nil
}

func (tc *testComposerConfig) GetXDebugCliIniPath() string {
	return ""
}

func (tc *testComposerConfig) GetXDebugFpmIniPath() string {
	return ""
}

func (tc *testComposerConfig) GetXDebugConifgLocaton() string {
	return ""
}

type testComposerDialog struct{}

func (d *testComposerDialog) SetMainContaner([]string) (int, string, error) {
	return 0, "", nil
}

func (d *testComposerDialog) StartCommand() (string, error) {
	return "", nil
}

func (d *testComposerDialog) StartDocker() (string, error) {
	return "", nil
}

func (d *testComposerDialog) SaveDockerProjectPath() (string, error) {
	return "", nil
}

func (d *testComposerDialog) DockerService() (string, error) {
	return "", nil
}

func (d *testComposerDialog) DockerProjectPath(p string) (string, error) {
	return "", nil
}

func (d *testComposerDialog) DockerCliXdebugIniFilePath(c string) (string, error) {
	return "", nil
}

func (d *testComposerDialog) DockerFpmXdebugIniFilePath(c string) (string, error) {
	return "", nil
}

func (d *testComposerDialog) XDebugConfigLocation() (int, string, error) {
	return 0, "", nil
}

type testComposer struct {
	containerList []string
	locaton       func(string, string) (string, error)
	ctype         string
	command       string
}

func (tc *testComposer) GetCommandLocaton() func(string, string) (string, error) {
	return tc.locaton
}

func (tc *testComposer) GetCallType() string {
	return tc.ctype
}

func (tc *testComposer) GetComposerCommand() string {
	return tc.command
}

func (tc *testComposer) GetContainerList() []string {
	return tc.containerList
}

func TestParseCommand(t *testing.T) {

	shortcommand, calltype, dockercmd := parseCommand("composer:update:memory")
	assert.EqualValues(t, []string{"update:memory", "memory", "update"}, []string{shortcommand, calltype, dockercmd})

	shortcommand, calltype, dockercmd = parseCommand("composer:update")
	assert.EqualValues(t, []string{"update", "", "update"}, []string{shortcommand, calltype, dockercmd})

	shortcommand, calltype, dockercmd = parseCommand("composer:install:memory")
	assert.EqualValues(t, []string{"install:memory", "memory", "install"}, []string{shortcommand, calltype, dockercmd})

	shortcommand, calltype, dockercmd = parseCommand("composer:install")
	assert.EqualValues(t, []string{"install", "", "install"}, []string{shortcommand, calltype, dockercmd})

	shortcommand, calltype, dockercmd = parseCommand("composer")
	assert.EqualValues(t, []string{"composer", "", ""}, []string{shortcommand, calltype, dockercmd})

	shortcommand, calltype, dockercmd = parseCommand("composer:memory")
	assert.EqualValues(t, []string{"composer:memory", "memory", ""}, []string{shortcommand, calltype, dockercmd})
}

func TestComposerHandleCase1(t *testing.T) {
	cfg := &testComposerConfig{
		mainContainer: "",
	}

	dlg := &testComposerDialog{}

	cmp := &testComposer{}

	a := &args{
		get:   "",
		slice: []string{},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{},
	}

	_, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.EqualError(t, err, "Container name is empty. Set the container name")
}

func TestComposerHandleCase2(t *testing.T) {
	cfg := &testComposerConfig{
		mainContainer: "containerName",
	}

	dlg := &testComposerDialog{}

	cmp := &testComposer{}

	a := &args{
		get:   "",
		slice: []string{},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{},
	}

	_, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.Nil(t, err)
}

func TestComposerHandleCase3(t *testing.T) {
	cfg := &testComposerConfig{
		mainContainer: "containerName",
	}

	dlg := &testComposerDialog{}

	cmp := &testComposer{
		locaton: func(container string, service string) (string, error) {
			return "/path/to/" + service, nil
		},
		ctype:   "",
		command: "",
	}

	a := &args{
		get:   "",
		slice: []string{},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{},
	}

	args, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.Nil(t, err)
	assert.Equal(t, []string{"exec", "-it", "containerName", "composer"}, args)
}

func TestComposerHandleCase4(t *testing.T) {
	cfg := &testComposerConfig{
		mainContainer: "containerName",
	}

	dlg := &testComposerDialog{}

	cmp := &testComposer{
		locaton: func(container string, service string) (string, error) {
			return "/path/to/" + service, nil
		},
		ctype:   "",
		command: "",
	}

	a := &args{
		get: "",
		slice: []string{
			"update",
		},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{},
	}

	args, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.Nil(t, err)
	assert.Equal(t, []string{"exec", "-it", "containerName", "composer", "update"}, args)
}

func TestComposerHandleCase5(t *testing.T) {
	cfg := &testComposerConfig{
		mainContainer: "containerName",
	}

	dlg := &testComposerDialog{}

	cmp := &testComposer{
		locaton: func(container string, service string) (string, error) {
			return "/path/to/" + service, nil
		},
		ctype:   "",
		command: "",
	}

	a := &args{
		get: "m",
		slice: []string{
			"m",
			"update",
		},
		tail: []string{
			"update",
		},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{},
	}

	args, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.Nil(t, err)
	assert.Equal(t, []string{"exec", "-i", "containerName", "/path/to/php", "-d", "memory_limit=-1", "/path/to/composer", "update"}, args)
}

func TestComposerHandleCase6(t *testing.T) {
	cfg := &testComposerConfig{
		mainContainer: "containerName",
	}

	dlg := &testComposerDialog{}

	cmp := &testComposer{
		locaton: func(container string, service string) (string, error) {
			return "/path/to/" + service, nil
		},
		ctype:   "memory",
		command: "update",
	}

	a := &args{
		get: "",
		slice: []string{
			"--help",
		},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{},
	}

	args, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.Nil(t, err)
	assert.Equal(t, []string{"exec", "-i", "containerName", "/path/to/php", "-d", "memory_limit=-1", "/path/to/composer", "update", "--help"}, args)
}

func TestComposerHandleCase7(t *testing.T) {
	cfg := &testComposerConfig{
		mainContainer: "containerName",
	}

	dlg := &testComposerDialog{}

	cmp := &testComposer{
		locaton: func(container string, service string) (string, error) {
			if service == "php" {
				return "", errors.New("Error on getting php path")
			}
			return "/path/to/" + service, nil
		},
		ctype:   "memory",
		command: "update",
	}

	a := &args{
		get: "",
		slice: []string{
			"--help",
		},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{},
	}

	_, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.EqualError(t, err, "Error on getting php path")
}

func TestComposerHandleCase8(t *testing.T) {
	cfg := &testComposerConfig{
		mainContainer: "containerName",
	}

	dlg := &testComposerDialog{}

	cmp := &testComposer{
		locaton: func(container string, service string) (string, error) {
			if service == "composer" {
				return "", errors.New("Error on getting composer path")
			}
			return "/path/to/" + service, nil
		},
		ctype:   "memory",
		command: "update",
	}

	a := &args{
		get: "",
		slice: []string{
			"--help",
		},
	}

	cl := &testContainerlist{
		err:           nil,
		containerList: []string{},
	}

	_, err := composerHandle(cfg, dlg, cmp, cl, a)

	assert.EqualError(t, err, "Error on getting composer path")
}
