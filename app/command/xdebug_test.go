package command

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testGetXdebugArgsProjectConfig struct {
	xDebugFpmIniPath     string
	xDebugCliIniPath     string
	xDebugConfigLocation string
	projectMainContainer string
}

func (c *testGetXdebugArgsProjectConfig) GetXDebugFpmIniPath() string {
	return c.xDebugFpmIniPath
}

func (c *testGetXdebugArgsProjectConfig) GetXDebugCliIniPath() string {
	return c.xDebugCliIniPath
}

func (c *testGetXdebugArgsProjectConfig) GetXDebugConfigLocaton() string {
	return c.xDebugConfigLocation
}

func (c *testGetXdebugArgsProjectConfig) GetProjectMainContainer() string {
	return c.projectMainContainer
}

func TestGetXdebugArgsCase1(t *testing.T) {
	cfg := &testGetXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugCliIniPath:     "/path/to/xdebug/cli.ini",
		xDebugConfigLocation: "container",
	}

	assert.EqualValues(t, []string{"docker", "exec", "main_container", "sed", "-i", "-e", `s/^\;zend_extension/zend_extension/g`, "/path/to/xdebug/cli.ini"}, getXdebugArgs(cfg, "xdebug:cli:enable", "/project/path/"))
}

func TestGetXdebugArgsCase2(t *testing.T) {
	cfg := &testGetXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugConfigLocation: "container",
		xDebugFpmIniPath:     "/path/to/xdebug/fpm.ini",
	}

	assert.EqualValues(t, []string{"docker", "exec", "main_container", "sed", "-i", "-e", `s/^\;zend_extension/zend_extension/g`, "/path/to/xdebug/fpm.ini"}, getXdebugArgs(cfg, "xdebug:fpm:enable", "/project/path/"))
}

func TestGetXdebugArgsCase3(t *testing.T) {
	cfg := &testGetXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugCliIniPath:     "/path/to/xdebug/cli.ini",
		xDebugConfigLocation: "container",
	}

	assert.EqualValues(t, []string{"docker", "exec", "main_container", "sed", "-i", "-e", `s/^zend_extension/\;zend_extension/g`, "/path/to/xdebug/cli.ini"}, getXdebugArgs(cfg, "xdebug:cli:disable", "/project/path/"))
}

func TestGetXdebugArgsCase4(t *testing.T) {
	cfg := &testGetXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugFpmIniPath:     "/path/to/xdebug/fpm.ini",
		xDebugConfigLocation: "container",
	}

	assert.EqualValues(t, []string{"docker", "exec", "main_container", "sed", "-i", "-e", `s/^zend_extension/\;zend_extension/g`, "/path/to/xdebug/fpm.ini"}, getXdebugArgs(cfg, "xdebug:fpm:disable", "/project/path/"))
}

func TestGetXdebugArgsCase5(t *testing.T) {
	cfg := &testGetXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugCliIniPath:     "/path/to/xdebug/cli.ini",
		xDebugConfigLocation: "local",
	}

	assert.EqualValues(t, []string{"sed", "-i", "-e", `s/^\;zend_extension/zend_extension/g`, "/project/path/path/to/xdebug/cli.ini"}, getXdebugArgs(cfg, "xdebug:cli:enable", "/project/path/"))
}

func TestGetXdebugArgsCase6(t *testing.T) {
	cfg := &testGetXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugConfigLocation: "local",
		xDebugFpmIniPath:     "/path/to/xdebug/fpm.ini",
	}

	assert.EqualValues(t, []string{"sed", "-i", "-e", `s/^\;zend_extension/zend_extension/g`, "/project/path/path/to/xdebug/fpm.ini"}, getXdebugArgs(cfg, "xdebug:fpm:enable", "/project/path"))
}

func TestGetXdebugArgsCase7(t *testing.T) {
	cfg := &testGetXdebugArgsProjectConfig{
		projectMainContainer: "main_container",
		xDebugCliIniPath:     "/path/to/xdebug/cli.ini",
		xDebugConfigLocation: "local",
	}

	assert.EqualValues(t, []string{"sed", "-i", "-e", `s/^zend_extension/\;zend_extension/g`, "/project/path/path/to/xdebug/cli.ini"}, getXdebugArgs(cfg, "xdebug:cli:disable", "/project/path/"))
}

func TestGetXdebugArgsCase8(t *testing.T) {
	cfg := &testGetXdebugArgsProjectConfig{
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

/*
type testXdebugProjectConfig struct{}

func (c *testXdebugProjectConfig) GetXDebugFpmIniPath() string {
	return ""
}

func (c *testXdebugProjectConfig) GetXDebugCliIniPath() string {
	return ""
}

func (c *testXdebugProjectConfig) GetXDebugConfigLocaton() string {
	return ""
}

func (c *testXdebugProjectConfig) GetProjectMainContainer() string {
	return ""
}

func (c *testXdebugProjectConfig) SaveContainerNameToProjectConfig(s string) error {
	return nil
}

func (c *testXdebugProjectConfig) SaveDockerCliXdebugIniFilePath(s string) error {
	return nil
}

func (c *testXdebugProjectConfig) SaveDockerFpmXdebugIniFilePath(s string) error {
	return nil
}

func (c *testXdebugProjectConfig) SaveXDebugConifgLocaton(s string) error {
	return nil
}
*/
