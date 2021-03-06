package command

import (
	"errors"
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

type syncConfig struct {
	projectMainContainer string
	projectDockerPath    string
}

func (sc *syncConfig) GetProjectMainContainer() string {
	return sc.projectMainContainer
}
func (sc *syncConfig) GetProjectDockerPath() string {
	return sc.projectDockerPath
}
func (sc *syncConfig) SaveContainerNameToProjectConfig(s string) error {
	return nil
}
func (sc *syncConfig) SaveStartCommandToProjectConfig(s string) error {
	return nil
}
func (sc *syncConfig) SaveDockerProjectPath(s string) error {
	return nil
}

type syncDlg struct {
	setMainContaner   func() (int, string, error)
	dockerProjectPath func() (string, error)
}

func (dlg *syncDlg) SetMainContaner([]string) (int, string, error) {
	return dlg.setMainContaner()
}
func (dlg *syncDlg) DockerProjectPath(string) (string, error) {
	return dlg.dockerProjectPath()
}

type testSyncOptions struct {
	getExecCommand   func(ExecOptions, *cli.App) error
	getInitFunction  func(bool) string
	getContainerList func() ([]string, error)
	getCopyTo        func(container string, sourcePath string, dstPath string) error
	runNativeExec    func(ExecOptions, *cli.App) error
	dirExists        func(string) (bool, error)
	mkdirAll         func(string, os.FileMode) error
}

func (x *testSyncOptions) GetExecCommand() func(ExecOptions, *cli.App) error {
	return x.getExecCommand
}
func (x *testSyncOptions) GetInitFunction() func(bool) string {
	return x.getInitFunction
}
func (x *testSyncOptions) GetContainerList() ([]string, error) {
	return x.getContainerList()
}
func (x *testSyncOptions) GetCopyTo(container string, sourcePath string, dstPath string) error {
	return x.getCopyTo(container, sourcePath, dstPath)
}
func (x *testSyncOptions) RunNativeExec(o ExecOptions, app *cli.App) error {
	return x.runNativeExec(o, app)
}
func (x *testSyncOptions) DirExists(path string) (bool, error) {
	return x.dirExists(path)
}
func (x *testSyncOptions) MkdirAll(path string, fileMode os.FileMode) error {
	return x.mkdirAll(path, fileMode)
}

func TestGetSyncPath(t *testing.T) {
	assert.Equal(t, "/vendor/path", getSyncPath("/vendor/path/"))
	assert.Equal(t, "/vendor/path", getSyncPath("/vendor/path"))
	assert.Equal(t, "/./", getSyncPath("./"))
	assert.Equal(t, "/./", getSyncPath("--all"))
}

func TestGetSyncArgs(t *testing.T) {
	cfg := &syncConfig{
		projectMainContainer: "phpfmp",
		projectDockerPath:    "/var/www/html/",
	}
	args := getSyncArgs(cfg, "copyfrom", getSyncPath("/vendor/path"), "/path/to/local/project/")
	assert.Equal(t, "cp phpfmp:/var/www/html/vendor/path /path/to/local/project/vendor", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyfrom", getSyncPath("/vendor/path/file.php"), "/path/to/local/project/")
	assert.Equal(t, "cp phpfmp:/var/www/html/vendor/path/file.php /path/to/local/project/vendor/path", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("/vendor/path/file.php"), "/path/to/local/project/")
	assert.Equal(t, "cp /path/to/local/project/vendor/path/file.php phpfmp:/var/www/html/vendor/path", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("/vendor/path/"), "/path/to/local/project/")
	assert.Equal(t, "cp /path/to/local/project/vendor/path phpfmp:/var/www/html/vendor", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("/vendor/path"), "/path/to/local/project/")
	assert.Equal(t, "cp /path/to/local/project/vendor/path phpfmp:/var/www/html/vendor", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("/vendor/path"), "/path/to/local/project")
	assert.Equal(t, "cp /path/to/local/project/vendor/path phpfmp:/var/www/html/vendor", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("--all"), "/path/to/local/project")
	assert.Equal(t, "cp /path/to/local/project/./ phpfmp:/var/www/html/", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyfrom", getSyncPath("--all"), "/path/to/local/project")
	assert.Equal(t, "cp phpfmp:/var/www/html/./ /path/to/local/project/", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyfrom", getSyncPath("/vendor/path"), "/path/to/local/project/")
	assert.Equal(t, "cp phpfmp:/var/www/html/vendor/path /path/to/local/project/vendor", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyfrom", getSyncPath("vendor/path/file.php"), "/path/to/local/project/")
	assert.Equal(t, "cp phpfmp:/var/www/html/vendor/path/file.php /path/to/local/project/vendor/path", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("vendor/path/file.php"), "/path/to/local/project/")
	assert.Equal(t, "cp /path/to/local/project/vendor/path/file.php phpfmp:/var/www/html/vendor/path", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("vendor/path/"), "/path/to/local/project/")
	assert.Equal(t, "cp /path/to/local/project/vendor/path phpfmp:/var/www/html/vendor", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("/vendor/path"), "/path/to/local/project/")
	assert.Equal(t, "cp /path/to/local/project/vendor/path phpfmp:/var/www/html/vendor", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("/vendor/path"), "/path/to/local/project")
	assert.Equal(t, "cp /path/to/local/project/vendor/path phpfmp:/var/www/html/vendor", strings.Join(args, " "))
}

func TestSyncCommandCase1(t *testing.T) {
	cfg := &syncConfig{}
	dlg := &syncDlg{}
	opt := &testSyncOptions{
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

	set := &flag.FlagSet{}
	set.Parse([]string{})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := SyncCommand("copyto", cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "Please, specify the path you want to sync")
}

func TestSyncCommandCase2(t *testing.T) {
	cfg := &syncConfig{}
	dlg := &syncDlg{}
	opt := &testSyncOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return nil
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, errors.New("get container list error")
		},
	}

	set := &flag.FlagSet{}
	set.Parse([]string{
		"path/to/sync",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := SyncCommand("copyto", cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "get container list error")
}

func TestSyncCommandCase3(t *testing.T) {
	cfg := &syncConfig{}
	dlg := &syncDlg{
		setMainContaner: func() (int, string, error) {
			return 0, "", errors.New("defineProjectMainContainer error")
		},
	}
	opt := &testSyncOptions{
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

	set := &flag.FlagSet{}
	set.Parse([]string{
		"path/to/sync",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := SyncCommand("copyto", cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "defineProjectMainContainer error")
}

func TestSyncCommandCase4(t *testing.T) {
	cfg := &syncConfig{}
	dlg := &syncDlg{
		setMainContaner: func() (int, string, error) {
			return 0, "main_container", nil
		},
		dockerProjectPath: func() (string, error) {
			return "", errors.New("DockerProjectPath error")
		},
	}
	opt := &testSyncOptions{
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

	set := &flag.FlagSet{}
	set.Parse([]string{
		"path/to/sync",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := SyncCommand("copyto", cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "DockerProjectPath error")
}

func TestSyncCommandCase5(t *testing.T) {
	cfg := &syncConfig{}
	dlg := &syncDlg{
		setMainContaner: func() (int, string, error) {
			return 0, "main_container", nil
		},
		dockerProjectPath: func() (string, error) {
			return "/var/www/html/", nil
		},
	}
	opt := &testSyncOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return errors.New("getExecCommand error")
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
	}

	set := &flag.FlagSet{}
	set.Bool("f", false, "")
	set.Parse([]string{
		"path/to/sync",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := SyncCommand("copyfrom", cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "getExecCommand error")
}

func TestSyncCommandCase6(t *testing.T) {
	cfg := &syncConfig{}
	dlg := &syncDlg{
		setMainContaner: func() (int, string, error) {
			return 0, "main_container", nil
		},
		dockerProjectPath: func() (string, error) {
			return "/var/www/html/", nil
		},
	}
	opt := &testSyncOptions{
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

	set := &flag.FlagSet{}
	set.Bool("f", false, "")
	set.Parse([]string{
		"path/to/sync",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := SyncCommand("copyfrom", cfg, dlg, opt)

	assert.Nil(t, app.Action(ctx))
}

func TestSyncCommandCase7(t *testing.T) {
	cfg := &syncConfig{}
	dlg := &syncDlg{
		setMainContaner: func() (int, string, error) {
			return 0, "main_container", nil
		},
		dockerProjectPath: func() (string, error) {
			return "/var/www/html/", nil
		},
	}
	opt := &testSyncOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return nil
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		mkdirAll: func(path string, filePerm os.FileMode) error {
			return errors.New("Path cannot be created")
		},
	}

	set := &flag.FlagSet{}
	set.Bool("f", true, "")
	set.Parse([]string{
		"path/to/sync",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := SyncCommand("copyfrom", cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "Path cannot be created")
}

func TestSyncCommandCase8(t *testing.T) {
	cfg := &syncConfig{}
	dlg := &syncDlg{
		setMainContaner: func() (int, string, error) {
			return 0, "main_container", nil
		},
		dockerProjectPath: func() (string, error) {
			return "/var/www/html/", nil
		},
	}
	opt := &testSyncOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return errors.New("ExecCommand error")
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		mkdirAll: func(path string, filePerm os.FileMode) error {
			return nil
		},
	}

	set := &flag.FlagSet{}
	set.Bool("f", true, "")
	set.Parse([]string{
		"path/to/sync",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := SyncCommand("copyfrom", cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "ExecCommand error")
}

func TestSyncCommandCase9(t *testing.T) {
	cfg := &syncConfig{}
	dlg := &syncDlg{
		setMainContaner: func() (int, string, error) {
			return 0, "main_container", nil
		},
		dockerProjectPath: func() (string, error) {
			return "/var/www/html/", nil
		},
	}
	opt := &testSyncOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return errors.New("ExecCommand error")
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		dirExists: func(path string) (bool, error) {
			return false, nil
		},
	}

	set := &flag.FlagSet{}
	set.Bool("f", true, "")
	set.Parse([]string{
		"path/to/sync",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := SyncCommand("copyto", cfg, dlg, opt)

	assert.Nil(t, app.Action(ctx))
}

func TestSyncCommandCase10(t *testing.T) {
	cfg := &syncConfig{}
	dlg := &syncDlg{
		setMainContaner: func() (int, string, error) {
			return 0, "main_container", nil
		},
		dockerProjectPath: func() (string, error) {
			return "/var/www/html/", nil
		},
	}
	opt := &testSyncOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return errors.New("ExecCommand error")
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		dirExists: func(path string) (bool, error) {
			return false, errors.New("Dir exists error")
		},
	}

	set := &flag.FlagSet{}
	set.Bool("f", true, "")
	set.Parse([]string{
		"path/to/sync",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := SyncCommand("copyto", cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "Dir exists error")
}

func TestSyncCommandCase11(t *testing.T) {
	cfg := &syncConfig{}
	dlg := &syncDlg{
		setMainContaner: func() (int, string, error) {
			return 0, "main_container", nil
		},
		dockerProjectPath: func() (string, error) {
			return "/var/www/html/", nil
		},
	}
	opt := &testSyncOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return errors.New("ExecCommand error")
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		dirExists: func(path string) (bool, error) {
			return true, nil
		},
		runNativeExec: func(o ExecOptions, ap *cli.App) error {
			return errors.New("path cannot be created")
		},
	}

	set := &flag.FlagSet{}
	set.Bool("f", true, "")
	set.Parse([]string{
		"path/to/sync",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := SyncCommand("copyto", cfg, dlg, opt)

	assert.EqualError(t, app.Action(ctx), "path cannot be created")
}

func TestSyncCommandCase12(t *testing.T) {
	cfg := &syncConfig{}
	dlg := &syncDlg{
		setMainContaner: func() (int, string, error) {
			return 0, "main_container", nil
		},
		dockerProjectPath: func() (string, error) {
			return "/var/www/html/", nil
		},
	}
	opt := &testSyncOptions{
		getExecCommand: func(o ExecOptions, a *cli.App) error {
			return nil
		},
		getInitFunction: func(s bool) string {
			return "/current/path"
		},
		getContainerList: func() ([]string, error) {
			return []string{"container"}, nil
		},
		dirExists: func(path string) (bool, error) {
			return true, nil
		},
		runNativeExec: func(o ExecOptions, ap *cli.App) error {
			return nil
		},
	}

	set := &flag.FlagSet{}
	set.Bool("f", true, "")
	set.Parse([]string{
		"path/to/sync",
	})

	ctx := &cli.Context{
		App: &cli.App{},
	}

	ctx = cli.NewContext(&cli.App{}, set, ctx)
	app := SyncCommand("copyto", cfg, dlg, opt)

	assert.Nil(t, app.Action(ctx))
}
