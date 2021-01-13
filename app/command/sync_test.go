package command

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func testGetSyncPath(t *testing.T) {
	assert.Equal(t, "/vendor/path", getSyncPath("/vendor/path/"))
	assert.Equal(t, "/vendor/path", getSyncPath("/vendor/path"))
	assert.Equal(t, "/./", getSyncPath("./"))
	assert.Equal(t, "/./", getSyncPath("--all"))
}

func testGetSyncArgs(t *testing.T) {
	cfg := &syncConfig{
		projectMainContainer: "phpfmp",
		projectDockerPath:    "/var/www/html/",
	}
	args := getSyncArgs(cfg, "copyfrom", getSyncPath("/vendor/path"), "/path/to/local/project/")
	assert.Equal(t, "cp phpfmp:/var/www/html/vendor/path /path/to/local/project/vendor/path", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyfrom", getSyncPath("/vendor/path/file.php"), "/path/to/local/project/")
	assert.Equal(t, "cp phpfmp:/var/www/html/vendor/path/file.php /path/to/local/project/vendor/path/file.php", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("/vendor/path/file.php"), "/path/to/local/project/")
	assert.Equal(t, "cp /path/to/local/project/vendor/path/file.php phpfmp:/var/www/html/vendor/path/file.php", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("/vendor/path/"), "/path/to/local/project/")
	assert.Equal(t, "cp /path/to/local/project/vendor/path phpfmp:/var/www/html/vendor/path", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("/vendor/path"), "/path/to/local/project/")
	assert.Equal(t, "cp /path/to/local/project/vendor/path phpfmp:/var/www/html/vendor/path", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("/vendor/path"), "/path/to/local/project")
	assert.Equal(t, "cp /path/to/local/project/vendor/path phpfmp:/var/www/html/vendor/path", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("--all"), "/path/to/local/project")
	assert.Equal(t, "cp /path/to/local/project/./ phpfmp:/var/www/html/./", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyfrom", getSyncPath("--all"), "/path/to/local/project")
	assert.Equal(t, "cp phpfmp:/var/www/html/./ /path/to/local/project/./", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyfrom", getSyncPath("/vendor/path"), "/path/to/local/project/")
	assert.Equal(t, "cp phpfmp:/var/www/html/vendor/path /path/to/local/project/vendor/path", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyfrom", getSyncPath("vendor/path/file.php"), "/path/to/local/project/")
	assert.Equal(t, "cp phpfmp:/var/www/html/vendor/path/file.php /path/to/local/project/vendor/path/file.php", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("vendor/path/file.php"), "/path/to/local/project/")
	assert.Equal(t, "cp /path/to/local/project/vendor/path/file.php phpfmp:/var/www/html/vendor/path/file.php", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("vendor/path/"), "/path/to/local/project/")
	assert.Equal(t, "cp /path/to/local/project/vendor/path phpfmp:/var/www/html/vendor/path", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("/vendor/path"), "/path/to/local/project/")
	assert.Equal(t, "cp /path/to/local/project/vendor/path phpfmp:/var/www/html/vendor/path", strings.Join(args, " "))

	args = getSyncArgs(cfg, "copyto", getSyncPath("/vendor/path"), "/path/to/local/project")
	assert.Equal(t, "cp /path/to/local/project/vendor/path phpfmp:/var/www/html/vendor/path", strings.Join(args, " "))
}
