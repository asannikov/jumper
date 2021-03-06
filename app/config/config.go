package config

import (
	"os"
	"strings"
)

type fileSystem interface {
	FileExists(string) (bool, error)
	DirExists(string) (bool, error)
	SaveConfigFile(interface{}, string) error
	ReadConfigFile(string, interface{}) error
	GoToProjectPath(string) error
}

// Config contains configuration from json config file
type Config struct {
	ProjectFile    string
	UserFile       string // has to contain absolute path to user config and the filename itself
	projectConfig  *ProjectConfig
	globalConfig   *GlobalConfig
	fileSystem     fileSystem
	hasProjectFile bool
}

// GetUserFile gets user config file
func (c *Config) GetUserFile() string {
	return c.UserFile
}

// SetUserFile sets user config file
func (c *Config) SetUserFile(file string) {
	c.UserFile = file
}

// ProjectSettings is not the same as ProjectConfig, but the similar one.
// it helps to call FindProjectPathInJSON from the outside (ie main function), where
// ProjectSettings is used as projectConfig, see main_start.go
type ProjectSettings interface {
	GetProjectName() string
	GetProjectPath() string
	SetProjectName(string)
	SetProjectPath(string)
}

// SetProjectPath set project path for config
func (c *Config) SetProjectPath(path string) {
	c.projectConfig.Path = path
}

// SetProjectName set project path for config
func (c *Config) SetProjectName(name string) {
	c.projectConfig.Name = name
}

// lookupProjectConfig seeks for a appropriate config
func (c *Config) lookupProjectConfig() (err error) {
	c.hasProjectFile = false
	if err = c.fileSystem.ReadConfigFile(c.ProjectFile, c.projectConfig); err == nil {
		c.hasProjectFile = true
		return nil
	}
	return err
}

// LookupUserConfig seeks for a user config
func (c *Config) lookupUserConfig() (err error) {
	if err = c.fileSystem.ReadConfigFile(c.GetUserFile(), c.globalConfig); err == nil {
		return nil
	}

	if err != nil && strings.Contains(err.Error(), "no such file or directory") == true {
		err = c.fileSystem.SaveConfigFile(c.globalConfig, c.GetUserFile())
	}

	return err
}

// Init Initiate conifg
func (c *Config) Init() {
	c.projectConfig = &ProjectConfig{}
	c.globalConfig = &GlobalConfig{}
}

// LoadProjectConfig loads project config
func (c *Config) LoadProjectConfig() (status bool, err error) {
	err = c.lookupProjectConfig()

	if err != nil && strings.Contains(err.Error(), "no such file or directory") == false {
		return false, err
	}

	return c.hasProjectFile, nil
}

// LoadConfig loads configuration
func (c *Config) LoadConfig(seekProject bool) (err error) {
	if err = c.lookupUserConfig(); err != nil {
		return err
	}

	if !seekProject {
		return nil
	}

	_, err = c.LoadProjectConfig()

	return err
}

// GetCommandInactveStatus gets command status
func (c *Config) GetCommandInactveStatus(cmd string) bool {
	return c.globalConfig.GetCommandInactveStatus(cmd)
}

// FindProjectPathInJSON check if project path in the json
func (c *Config) FindProjectPathInJSON(pc ProjectSettings) {
	c.globalConfig.FindProjectPathInJSON(func(p GlobalProjectConfig) (bool, error) {
		if p.Name == pc.GetProjectName() {
			if t, err := c.fileSystem.DirExists(p.Path); err == nil && t == true {
				pc.SetProjectName(p.Name)
				pc.SetProjectPath(p.Path)
				return true, nil
			} else if err != nil {
				return false, err
			}
		}

		return false, nil
	})
}

// GetProjectPath gets project path
func (c *Config) GetProjectPath() string {
	return c.projectConfig.GetPath()
}

// GetProjectName gets project path
func (c *Config) GetProjectName() string {
	return c.projectConfig.GetName()
}

// ProjectConfgFileFound checks if current path has config file
func (c *Config) ProjectConfgFileFound() bool {
	return c.hasProjectFile
}

// GetProjectNameList get projects name list
func (c *Config) GetProjectNameList() []string {
	return c.globalConfig.GetProjectNameList()
}

// SetFileSystem set file system object
func (c *Config) SetFileSystem(fs fileSystem) {
	c.fileSystem = fs
}

// AddProjectConfigFile generates project config file
func (c *Config) AddProjectConfigFile() (err error) {
	projectFile := strings.TrimRight(c.GetProjectPath(), string(os.PathSeparator)) + string(os.PathSeparator) + c.ProjectFile
	if err = c.fileSystem.SaveConfigFile(c.projectConfig, projectFile); err != nil {
		return err
	}

	fpc := GlobalProjectConfig{
		Name: c.GetProjectName(),
		Path: c.GetProjectPath(),
	}

	c.globalConfig.Projects = append(c.globalConfig.Projects, fpc)
	c.globalConfig.InactiveCommandTypes = []string{"composer", "php", "magento"}
	return c.fileSystem.SaveConfigFile(c.globalConfig, c.GetUserFile())
}

// GetFile gets project file
func (c *Config) GetFile() string {
	return c.ProjectFile
}

func (c *Config) saveProjectFile() error {
	return c.fileSystem.SaveConfigFile(c.projectConfig, c.getProjectFile())
}

func (c *Config) getProjectFile() string {
	return strings.TrimRight(c.projectConfig.GetPath(), string(os.PathSeparator)) + string(os.PathSeparator) + c.ProjectFile
}
