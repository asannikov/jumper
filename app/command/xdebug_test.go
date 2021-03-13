package command

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

type testXdebugArgsProjectConfig struct {
	xDebugFpmIniPath     string
	xDebugCliIniPath     string
	xDebugConfigLocation string
	projectMainContainer string

	saveContainerNameToProjectConfig error
	saveDockerCliXdebugIniFilePath   error
	saveDockerFpmXdebugIniFilePath   error
	saveXDebugConifgLocaton          error
	getCommandInactveStatus          bool
}

func (c *testXdebugArgsProjectConfig) GetXDebugFpmIniPath() string {
	return c.xDebugFpmIniPath
}
func (c *testXdebugArgsProjectConfig) GetXDebugCliIniPath() string {
	return c.xDebugCliIniPath
}
func (c *testXdebugArgsProjectConfig) GetXDebugConfigLocaton() string {
	return c.xDebugConfigLocation
}
func (c *testXdebugArgsProjectConfig) GetProjectMainContainer() string {
	return c.projectMainContainer
}
func (c *testXdebugArgsProjectConfig) SaveContainerNameToProjectConfig(v string) error {
	return c.saveContainerNameToProjectConfig
}
func (c *testXdebugArgsProjectConfig) SaveDockerCliXdebugIniFilePath(v string) error {
	return c.saveDockerCliXdebugIniFilePath
}
func (c *testXdebugArgsProjectConfig) SaveDockerFpmXdebugIniFilePath(v string) error {
	return c.saveDockerFpmXdebugIniFilePath
}
func (c *testXdebugArgsProjectConfig) SaveXDebugConifgLocaton(v string) error {
	return c.saveXDebugConifgLocaton
}
func (c *testXdebugArgsProjectConfig) GetCommandInactveStatus(v string) bool {
	return c.getCommandInactveStatus
}

type testXDebugCommandDialog struct {
	setMainContaner            func([]string) (int, string, error)
	dockerCliXdebugIniFilePath func(string) (string, error)
	dockerFpmXdebugIniFilePath func(string) (string, error)
	xDebugConfigLocation       func() (int, string, error)
}

func (x *testXDebugCommandDialog) SetMainContaner(list []string) (int, string, error) {
	return x.setMainContaner(list)
}
func (x *testXDebugCommandDialog) DockerCliXdebugIniFilePath(v string) (string, error) {
	return x.dockerCliXdebugIniFilePath(v)
}
func (x *testXDebugCommandDialog) DockerFpmXdebugIniFilePath(v string) (string, error) {
	return x.dockerFpmXdebugIniFilePath(v)
}
func (x *testXDebugCommandDialog) XDebugConfigLocation() (int, string, error) {
	return x.xDebugConfigLocation()
}

type testXDebugOptions struct {
	getExecCommand    func(ExecOptions, *cli.App) error
	getInitFunction   func(bool) string
	getContainerList  func() ([]string, error)
	checkXdebugStatus func(*cli.App, []string) (bool, error)
}

func (x *testXDebugOptions) GetExecCommand() func(ExecOptions, *cli.App) error {
	return x.getExecCommand
}
func (x *testXDebugOptions) GetInitFunction() func(bool) string {
	return x.getInitFunction
}
func (x *testXDebugOptions) GetContainerList() ([]string, error) {
	return x.getContainerList()
}
func (x *testXDebugOptions) CheckXdebugStatus(app *cli.App, args []string) (bool, error) {
	return x.checkXdebugStatus(app, args)
}

func TestToggleXdebugArgsCase1(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugFpmIniPath:     "/path/to/xdebug/fpm.ini",
		xDebugConfigLocation: "container",
	}

	assert.EqualValues(t, []string{"docker", "exec", "-u", "root", "main_container", "cat", "/path/to/xdebug/fpm.ini"}, toggleXdebugArgs(cfg, "/project/path/", "xdebug:fpm:toggle"))
}

func TestToggleXdebugArgsCase2(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugFpmIniPath:     "/path/to/xdebug/fpm.ini",
		xDebugCliIniPath:     "/path/to/xdebug/cli.ini",
		xDebugConfigLocation: "container",
	}

	assert.EqualValues(t, []string{"docker", "exec", "-u", "root", "main_container", "cat", "/path/to/xdebug/cli.ini"}, toggleXdebugArgs(cfg, "/project/path/", "xdebug:cli:toggle"))
}

func TestToggleXdebugArgsCase3(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugFpmIniPath:     "/path/to/xdebug/fpm.ini",
		xDebugCliIniPath:     "/path/to/xdebug/cli.ini",
		xDebugConfigLocation: "local",
	}

	assert.EqualValues(t, []string{"cat", "/project/path/path/to/xdebug/cli.ini"}, toggleXdebugArgs(cfg, "/project/path/", "xdebug:cli:toggle"))
}

func TestToggleXdebugArgsCase4(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugFpmIniPath:     "/path/to/xdebug/fpm.ini",
		xDebugCliIniPath:     "/path/to/xdebug/cli.ini",
		xDebugConfigLocation: "local",
	}

	assert.EqualValues(t, []string{"cat", "/project/path/path/to/xdebug/cli.ini"}, toggleXdebugArgs(cfg, "/project/path", "xdebug:cli:toggle"))
}

func TestGetXdebugArgsCase1(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugCliIniPath:     "/path/to/xdebug/cli.ini",
		xDebugConfigLocation: "container",
	}

	assert.EqualValues(t, []string{"docker", "exec", "-u", "root", "main_container", "sed", "-i", "-e", `s/^\;zend_extension/zend_extension/g`, "/path/to/xdebug/cli.ini"}, getXdebugArgs(cfg, "xdebug:cli:enable", "/project/path/"))
}

func TestGetXdebugArgsCase2(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugConfigLocation: "container",
		xDebugFpmIniPath:     "/path/to/xdebug/fpm.ini",
	}

	assert.EqualValues(t, []string{"docker", "exec", "-u", "root", "main_container", "sed", "-i", "-e", `s/^\;zend_extension/zend_extension/g`, "/path/to/xdebug/fpm.ini"}, getXdebugArgs(cfg, "xdebug:fpm:enable", "/project/path/"))
}

func TestGetXdebugArgsCase3(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugCliIniPath:     "/path/to/xdebug/cli.ini",
		xDebugConfigLocation: "container",
	}

	assert.EqualValues(t, []string{"docker", "exec", "-u", "root", "main_container", "sed", "-i", "-e", `s/^zend_extension/\;zend_extension/g`, "/path/to/xdebug/cli.ini"}, getXdebugArgs(cfg, "xdebug:cli:disable", "/project/path/"))
}

func TestGetXdebugArgsCase4(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugFpmIniPath:     "/path/to/xdebug/fpm.ini",
		xDebugConfigLocation: "container",
	}

	assert.EqualValues(t, []string{"docker", "exec", "-u", "root", "main_container", "sed", "-i", "-e", `s/^zend_extension/\;zend_extension/g`, "/path/to/xdebug/fpm.ini"}, getXdebugArgs(cfg, "xdebug:fpm:disable", "/project/path/"))
}

func TestGetXdebugArgsCase5(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugCliIniPath:     "/path/to/xdebug/cli.ini",
		xDebugConfigLocation: "local",
	}

	assert.EqualValues(t, []string{"sed", "-i", "-e", `s/^\;zend_extension/zend_extension/g`, "/project/path/path/to/xdebug/cli.ini"}, getXdebugArgs(cfg, "xdebug:cli:enable", "/project/path/"))
}

func TestGetXdebugArgsCase6(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugConfigLocation: "local",
		xDebugFpmIniPath:     "/path/to/xdebug/fpm.ini",
	}

	assert.EqualValues(t, []string{"sed", "-i", "-e", `s/^\;zend_extension/zend_extension/g`, "/project/path/path/to/xdebug/fpm.ini"}, getXdebugArgs(cfg, "xdebug:fpm:enable", "/project/path"))
}

func TestGetXdebugArgsCase7(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugCliIniPath:     "/path/to/xdebug/cli.ini",
		xDebugConfigLocation: "local",
	}

	assert.EqualValues(t, []string{"sed", "-i", "-e", `s/^zend_extension/\;zend_extension/g`, "/project/path/path/to/xdebug/cli.ini"}, getXdebugArgs(cfg, "xdebug:cli:disable", "/project/path/"))
}

func TestGetXdebugArgsCase8(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugFpmIniPath:     "/path/to/xdebug/fpm.ini",
		xDebugConfigLocation: "local",
	}

	assert.EqualValues(t, []string{"sed", "-i", "-e", `s/^zend_extension/\;zend_extension/g`, "/project/path/path/to/xdebug/fpm.ini"}, getXdebugArgs(cfg, "xdebug:fpm:disable", "/project/path/"))
}

type testDefineCliXdebugIniFilePathProjectConfig struct {
	getXDebugConfigLocaton  string
	saveXDebugConifgLocaton error
}

func (c *testDefineCliXdebugIniFilePathProjectConfig) GetXDebugCliIniPath() string {
	return c.getXDebugConfigLocaton
}

func (c *testDefineCliXdebugIniFilePathProjectConfig) SaveDockerCliXdebugIniFilePath(s string) error {
	return c.saveXDebugConifgLocaton
}

type testDefineCliXdebugIniFilePathDialog struct {
	xDebugConfigLocationPath string
	xDebugConfigLocationErr  error
}

func (d *testDefineCliXdebugIniFilePathDialog) DockerCliXdebugIniFilePath(s string) (string, error) {
	return d.xDebugConfigLocationPath, d.xDebugConfigLocationErr
}

func TestDefineCliXdebugIniFilePathCase1(t *testing.T) {
	cfg := &testDefineCliXdebugIniFilePathProjectConfig{}

	d := &testDefineCliXdebugIniFilePathDialog{
		xDebugConfigLocationPath: "/path/to/cli.ini",
		xDebugConfigLocationErr:  nil,
	}

	assert.Nil(t, defineCliXdebugIniFilePath(cfg, d, "/etc/php/7.0/cli/conf.d/xdebug.ini"))
}

func TestDefineCliXdebugIniFilePathCase2(t *testing.T) {
	cfg := &testDefineCliXdebugIniFilePathProjectConfig{}

	d := &testDefineCliXdebugIniFilePathDialog{
		xDebugConfigLocationPath: "",
		xDebugConfigLocationErr:  nil,
	}

	assert.EqualError(t, defineCliXdebugIniFilePath(cfg, d, "/etc/php/7.0/cli/conf.d/xdebug.ini"), "Cli Xdebug ini file path is empty")
}

func TestDefineCliXdebugIniFilePathCase3(t *testing.T) {
	cfg := &testDefineCliXdebugIniFilePathProjectConfig{}

	d := &testDefineCliXdebugIniFilePathDialog{
		xDebugConfigLocationPath: "",
		xDebugConfigLocationErr:  errors.New("Dialog problem"),
	}

	assert.EqualError(t, defineCliXdebugIniFilePath(cfg, d, "/etc/php/7.0/cli/conf.d/xdebug.ini"), "Dialog problem")
}

func TestDefineCliXdebugIniFilePathCase4(t *testing.T) {
	cfg := &testDefineCliXdebugIniFilePathProjectConfig{
		saveXDebugConifgLocaton: errors.New("Error on saving"),
	}

	d := &testDefineCliXdebugIniFilePathDialog{
		xDebugConfigLocationPath: "/path/to/fpm.ini",
		xDebugConfigLocationErr:  nil,
	}

	assert.EqualError(t, defineCliXdebugIniFilePath(cfg, d, "/etc/php/7.0/cli/conf.d/xdebug.ini"), "Error on saving")
}

func TestDefineCliXdebugIniFilePathCase5(t *testing.T) {
	cfg := &testDefineCliXdebugIniFilePathProjectConfig{
		getXDebugConfigLocaton: "/path/to/cli.ini",
	}

	d := &testDefineCliXdebugIniFilePathDialog{
		xDebugConfigLocationPath: "/path/to/cli.ini",
		xDebugConfigLocationErr:  errors.New("Dialog problem"),
	}

	assert.Nil(t, defineCliXdebugIniFilePath(cfg, d, "/etc/php/7.0/cli/conf.d/xdebug.ini"))
}

type testDefineXdebugIniFileLocationProjectConfig struct {
	getXDebugConfigLocaton  string
	saveXDebugConifgLocaton error
}

func (c *testDefineXdebugIniFileLocationProjectConfig) SaveXDebugConifgLocaton(s string) error {
	return c.saveXDebugConifgLocaton
}

func (c *testDefineXdebugIniFileLocationProjectConfig) GetXDebugConfigLocaton() string {
	return c.getXDebugConfigLocaton
}

type testDefineXdebugIniFileLocationDialog struct {
	xDebugConfigLocationInt  int
	xDebugConfigLocationPath string
	xDebugConfigLocationErr  error
}

func (d *testDefineXdebugIniFileLocationDialog) XDebugConfigLocation() (int, string, error) {
	return d.xDebugConfigLocationInt, d.xDebugConfigLocationPath, d.xDebugConfigLocationErr
}

func TestDefineXdebugIniFileLocationCase1(t *testing.T) {
	cfg := &testDefineXdebugIniFileLocationProjectConfig{}

	d := &testDefineXdebugIniFileLocationDialog{
		xDebugConfigLocationInt:  0,
		xDebugConfigLocationPath: "/path/to/fpm.ini",
		xDebugConfigLocationErr:  nil,
	}

	assert.Nil(t, defineXdebugIniFileLocation(cfg, d))
}

func TestDefineXdebugIniFileLocationCase2(t *testing.T) {
	cfg := &testDefineXdebugIniFileLocationProjectConfig{}

	d := &testDefineXdebugIniFileLocationDialog{
		xDebugConfigLocationInt:  0,
		xDebugConfigLocationPath: "",
		xDebugConfigLocationErr:  nil,
	}

	assert.EqualError(t, defineXdebugIniFileLocation(cfg, d), "Xdebug config file locaton cannot be empty")
}

func TestDefineXdebugIniFileLocationCase3(t *testing.T) {
	cfg := &testDefineXdebugIniFileLocationProjectConfig{}

	d := &testDefineXdebugIniFileLocationDialog{
		xDebugConfigLocationInt:  0,
		xDebugConfigLocationPath: "",
		xDebugConfigLocationErr:  errors.New("Dialog problem"),
	}

	assert.EqualError(t, defineXdebugIniFileLocation(cfg, d), "Dialog problem")
}

func TestDefineXdebugIniFileLocationCase4(t *testing.T) {
	cfg := &testDefineXdebugIniFileLocationProjectConfig{
		saveXDebugConifgLocaton: errors.New("Error on saving"),
	}

	d := &testDefineXdebugIniFileLocationDialog{
		xDebugConfigLocationInt:  1,
		xDebugConfigLocationPath: "/path/to/fpm.ini",
		xDebugConfigLocationErr:  nil,
	}

	assert.EqualError(t, defineXdebugIniFileLocation(cfg, d), "Error on saving")
}

func TestDefineXdebugIniFileLocationCase5(t *testing.T) {
	cfg := &testDefineXdebugIniFileLocationProjectConfig{
		getXDebugConfigLocaton: "/path/to/fpm.ini",
	}

	d := &testDefineXdebugIniFileLocationDialog{
		xDebugConfigLocationInt:  0,
		xDebugConfigLocationPath: "",
		xDebugConfigLocationErr:  errors.New("Dialog problem"),
	}

	assert.Nil(t, defineXdebugIniFileLocation(cfg, d))
}

// defineFpmXdebugIniFilePath

type testDefineFpmXdebugIniFilePathProjectConfig struct {
	getXDebugConfigLocaton  string
	saveXDebugConifgLocaton error
}

func (c *testDefineFpmXdebugIniFilePathProjectConfig) GetXDebugFpmIniPath() string {
	return c.getXDebugConfigLocaton
}

func (c *testDefineFpmXdebugIniFilePathProjectConfig) SaveDockerFpmXdebugIniFilePath(s string) error {
	return c.saveXDebugConifgLocaton
}

type testDefineFpmXdebugIniFilePathDialog struct {
	xDebugConfigLocationPath string
	xDebugConfigLocationErr  error
}

func (d *testDefineFpmXdebugIniFilePathDialog) DockerFpmXdebugIniFilePath(s string) (string, error) {
	return d.xDebugConfigLocationPath, d.xDebugConfigLocationErr
}

func TestDefineFpmXdebugIniFilePathCase1(t *testing.T) {
	cfg := &testDefineFpmXdebugIniFilePathProjectConfig{}

	d := &testDefineFpmXdebugIniFilePathDialog{
		xDebugConfigLocationPath: "/path/to/fpm.ini",
		xDebugConfigLocationErr:  nil,
	}

	assert.Nil(t, defineFpmXdebugIniFilePath(cfg, d, "/etc/php/7.0/cli/conf.d/xdebug.ini"))
}

func TestDefineFpmXdebugIniFilePathCase2(t *testing.T) {
	cfg := &testDefineFpmXdebugIniFilePathProjectConfig{}

	d := &testDefineFpmXdebugIniFilePathDialog{
		xDebugConfigLocationPath: "",
		xDebugConfigLocationErr:  nil,
	}

	assert.EqualError(t, defineFpmXdebugIniFilePath(cfg, d, "/etc/php/7.0/cli/conf.d/xdebug.ini"), "Fpm Xdebug ini file path is empty")
}

func TestDefineFpmXdebugIniFilePathCase3(t *testing.T) {
	cfg := &testDefineFpmXdebugIniFilePathProjectConfig{}

	d := &testDefineFpmXdebugIniFilePathDialog{
		xDebugConfigLocationPath: "",
		xDebugConfigLocationErr:  errors.New("Dialog problem"),
	}

	assert.EqualError(t, defineFpmXdebugIniFilePath(cfg, d, "/etc/php/7.0/fpm/conf.d/xdebug.ini"), "Dialog problem")
}

func TestDefineFpmXdebugIniFilePathCase4(t *testing.T) {
	cfg := &testDefineFpmXdebugIniFilePathProjectConfig{
		saveXDebugConifgLocaton: errors.New("Error on saving"),
	}

	d := &testDefineFpmXdebugIniFilePathDialog{
		xDebugConfigLocationPath: "/path/to/fpm.ini",
		xDebugConfigLocationErr:  nil,
	}

	assert.EqualError(t, defineFpmXdebugIniFilePath(cfg, d, "/etc/php/7.0/fpm/conf.d/xdebug.ini"), "Error on saving")
}

func TestDefineFpmXdebugIniFilePathCase5(t *testing.T) {
	cfg := &testDefineFpmXdebugIniFilePathProjectConfig{
		getXDebugConfigLocaton: "/path/to/fpm.ini",
	}

	d := &testDefineFpmXdebugIniFilePathDialog{
		xDebugConfigLocationPath: "/path/to/fpm.ini",
		xDebugConfigLocationErr:  errors.New("Dialog problem"),
	}

	assert.Nil(t, defineFpmXdebugIniFilePath(cfg, d, "/etc/php/7.0/fpm/conf.d/xdebug.ini"))
}

func TestXDebugCommandXdebugFpmToggleCase1(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugFpmIniPath:     "/path/to/xdebug/fpm.ini",
		xDebugCliIniPath:     "/path/to/xdebug/cli.ini",
		xDebugConfigLocation: "container",
	}
	dlg := &testXDebugCommandDialog{}
	opt := &testXDebugOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return nil
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		checkXdebugStatus: func(a *cli.App, args []string) (bool, error) {
			return false, errors.New("CheckXdebugStatus error")
		},
	}
	app := XDebugCommand("xdebug:fpm:toggle", cfg, dlg, opt)
	ctx := &cli.Context{
		App: &cli.App{},
	}

	assert.EqualError(t, app.Action(ctx), "CheckXdebugStatus error")
}

func TestXDebugCommandXdebugFpmToggleCase2(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugFpmIniPath:     "/path/to/xdebug/fpm.ini",
		xDebugCliIniPath:     "/path/to/xdebug/cli.ini",
		xDebugConfigLocation: "container",
	}
	dlg := &testXDebugCommandDialog{}
	opt := &testXDebugOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			assert.EqualValues(t, []string{"exec", "-u", "root", "main_container", "sed", "-i", "-e", "s/^\\;zend_extension/zend_extension/g", "/path/to/xdebug/fpm.ini"}, o.GetArgs())
			return errors.New("Stop exec command")
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		checkXdebugStatus: func(a *cli.App, args []string) (bool, error) {
			return false, nil
		},
	}
	app := XDebugCommand("xdebug:fpm:toggle", cfg, dlg, opt)
	ctx := &cli.Context{
		App: &cli.App{},
	}
	app.Action(ctx)
}

func TestXDebugCommandXdebugFpmToggleCase3(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugFpmIniPath:     "/path/to/xdebug/fpm.ini",
		xDebugCliIniPath:     "/path/to/xdebug/cli.ini",
		xDebugConfigLocation: "container",
	}
	dlg := &testXDebugCommandDialog{}
	opt := &testXDebugOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			assert.EqualValues(t, []string{"exec", "-u", "root", "main_container", "sed", "-i", "-e", "s/^zend_extension/\\;zend_extension/g", "/path/to/xdebug/fpm.ini"}, o.GetArgs())
			return errors.New("Stop exec command")
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		checkXdebugStatus: func(a *cli.App, args []string) (bool, error) {
			return true, nil
		},
	}
	app := XDebugCommand("xdebug:fpm:toggle", cfg, dlg, opt)
	ctx := &cli.Context{
		App: &cli.App{},
	}
	app.Action(ctx)
}

func TestXDebugCommandXdebugFpmEnableCase1(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{}
	dlg := &testXDebugCommandDialog{}
	opt := &testXDebugOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return nil
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, errors.New("getContainerList error")
		},
	}
	app := XDebugCommand("xdebug:fpm:enable", cfg, dlg, opt)
	ctx := &cli.Context{
		App: &cli.App{},
	}

	assert.EqualError(t, app.Action(ctx), "getContainerList error")
}

func TestXDebugCommandXdebugFpmEnableCase2(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{}
	dlg := &testXDebugCommandDialog{
		setMainContaner: func(list []string) (int, string, error) {
			return 0, "", errors.New("defineProjectMainContainer error")
		},
	}
	opt := &testXDebugOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return nil
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
	}
	app := XDebugCommand("xdebug:fpm:enable", cfg, dlg, opt)
	ctx := &cli.Context{
		App: &cli.App{},
	}

	assert.EqualError(t, app.Action(ctx), "defineProjectMainContainer error")
}

func TestXDebugCommandXdebugFpmEnableCase3(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		xDebugFpmIniPath:               "",
		saveDockerCliXdebugIniFilePath: errors.New("defineCliXdebugIniFilePath - SaveDockerCliXdebugIniFilePath error"),
		saveDockerFpmXdebugIniFilePath: errors.New("defineFpmXdebugIniFilePath - SaveDockerFpmXdebugIniFilePath error"),
	}
	dlg := &testXDebugCommandDialog{
		setMainContaner: func(list []string) (int, string, error) {
			return 1, "container_name", nil
		},
		dockerCliXdebugIniFilePath: func(path string) (string, error) {
			return "/path/to/cli", nil
		},
		dockerFpmXdebugIniFilePath: func(path string) (string, error) {
			return "/path/to/fpm", nil
		},
	}
	opt := &testXDebugOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return nil
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
	}
	app := XDebugCommand("xdebug:fpm:enable", cfg, dlg, opt)
	ctx := &cli.Context{
		App: &cli.App{},
	}

	assert.EqualError(t, app.Action(ctx), "defineFpmXdebugIniFilePath - SaveDockerFpmXdebugIniFilePath error")
}

func TestXDebugCommandXdebugCliEnableCase4(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		xDebugFpmIniPath:               "",
		saveDockerCliXdebugIniFilePath: errors.New("defineCliXdebugIniFilePath - SaveDockerCliXdebugIniFilePath error"),
		saveDockerFpmXdebugIniFilePath: errors.New("defineFpmXdebugIniFilePath - SaveDockerFpmXdebugIniFilePath error"),
	}
	dlg := &testXDebugCommandDialog{
		setMainContaner: func(list []string) (int, string, error) {
			return 1, "container_name", nil
		},
		dockerCliXdebugIniFilePath: func(path string) (string, error) {
			return "/path/to/cli", nil
		},
		dockerFpmXdebugIniFilePath: func(path string) (string, error) {
			return "/path/to/fpm", nil
		},
	}
	opt := &testXDebugOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return nil
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
	}
	app := XDebugCommand("xdebug:cli:enable", cfg, dlg, opt)
	ctx := &cli.Context{
		App: &cli.App{},
	}

	assert.EqualError(t, app.Action(ctx), "defineCliXdebugIniFilePath - SaveDockerCliXdebugIniFilePath error")
}

func TestXDebugCommandXdebugCliEnableCase5(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		xDebugFpmIniPath:               "",
		saveDockerCliXdebugIniFilePath: nil,
		saveDockerFpmXdebugIniFilePath: nil,
	}
	dlg := &testXDebugCommandDialog{
		setMainContaner: func(list []string) (int, string, error) {
			return 1, "container_name", nil
		},
		dockerCliXdebugIniFilePath: func(path string) (string, error) {
			return "/path/to/cli", nil
		},
		dockerFpmXdebugIniFilePath: func(path string) (string, error) {
			return "/path/to/fpm", nil
		},
		xDebugConfigLocation: func() (int, string, error) {
			return 0, "", errors.New("defineXdebugIniFileLocation error")
		},
	}
	opt := &testXDebugOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return nil
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
	}
	app := XDebugCommand("xdebug:cli:enable", cfg, dlg, opt)
	ctx := &cli.Context{
		App: &cli.App{},
	}

	assert.EqualError(t, app.Action(ctx), "defineXdebugIniFileLocation error")
}

func TestXDebugCommandXdebugCliEnableCase6(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		xDebugFpmIniPath:               "",
		saveDockerCliXdebugIniFilePath: nil,
		saveDockerFpmXdebugIniFilePath: nil,
	}
	dlg := &testXDebugCommandDialog{
		setMainContaner: func(list []string) (int, string, error) {
			return 1, "container_name", nil
		},
		dockerCliXdebugIniFilePath: func(path string) (string, error) {
			return "/path/to/cli", nil
		},
		dockerFpmXdebugIniFilePath: func(path string) (string, error) {
			return "/path/to/fpm", nil
		},
		xDebugConfigLocation: func() (int, string, error) {
			return 0, "local", nil
		},
	}
	opt := &testXDebugOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return errors.New("exec command error")
		},
	}
	app := XDebugCommand("xdebug:cli:enable", cfg, dlg, opt)
	ctx := &cli.Context{
		App: &cli.App{},
	}

	assert.EqualError(t, app.Action(ctx), "exec command error")
}

func TestXDebugCommandXdebug7(t *testing.T) {
	cfg := &testXdebugArgsProjectConfig{
		xDebugFpmIniPath:               "",
		saveDockerCliXdebugIniFilePath: nil,
		saveDockerFpmXdebugIniFilePath: nil,
	}
	dlg := &testXDebugCommandDialog{
		setMainContaner: func(list []string) (int, string, error) {
			return 1, "container_name", nil
		},
		dockerCliXdebugIniFilePath: func(path string) (string, error) {
			return "/path/to/cli", nil
		},
		dockerFpmXdebugIniFilePath: func(path string) (string, error) {
			return "/path/to/fpm", nil
		},
		xDebugConfigLocation: func() (int, string, error) {
			return 0, "local", nil
		},
	}
	opt := &testXDebugOptions{
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return nil
		},
	}

	ctx := &cli.Context{
		App: &cli.App{},
	}

	app := XDebugCommand("xdebug:cli:enable", cfg, dlg, opt)
	assert.Nil(t, app.Action(ctx))

	app = XDebugCommand("xdebug:cli:disable", cfg, dlg, opt)
	assert.Nil(t, app.Action(ctx))

	app = XDebugCommand("xdebug:fpm:enable", cfg, dlg, opt)
	assert.Nil(t, app.Action(ctx))

	app = XDebugCommand("xdebug:fpm:disable", cfg, dlg, opt)
	assert.Nil(t, app.Action(ctx))
}
