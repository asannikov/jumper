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

type projectSettings interface {
	GetProjectName() string
	GetProjectPath() string
	SetProjectName(string)
	SetProjectPath(string)
}

// LookupProjectConfig seeks for a appropriate config
func (c *Config) lookupProjectConfig() (err error) {
	ps := ProjectConfig{}
	c.hasProjectFile = false
	if err = c.fileSystem.ReadConfigFile(c.ProjectFile, &ps); err == nil {
		c.projectConfig = &ps
		c.hasProjectFile = true
		return nil
	}
	return err
}

// LookupUserConfig seeks for a user config
func (c *Config) lookupUserConfig() (err error) {
	gc := GlobalConfig{}
	c.globalConfig = &gc
	if err = c.fileSystem.ReadConfigFile(c.UserFile, &gc); err == nil {
		return nil
	}

	if err != nil && strings.Contains(err.Error(), "no such file or directory") == true {
		err = c.fileSystem.SaveConfigFile(gc, c.UserFile)
	}

	return err
}

// AddProjectConfigFile generates project config file
func (c *Config) AddProjectConfigFile(pc *ProjectConfig) (err error) {
	projectFile := strings.TrimRight(pc.GetPath(), string(os.PathSeparator)) + string(os.PathSeparator) + c.ProjectFile
	if err = c.fileSystem.SaveConfigFile(pc, projectFile); err != nil {
		return err
	}

	fpc := GlobalProjectConfig{
		Name: pc.GetName(),
		Path: pc.GetPath(),
	}

	c.globalConfig.Projects = append(c.globalConfig.Projects, fpc)

	return c.fileSystem.SaveConfigFile(c.globalConfig, c.UserFile)
}

// LoadConfig loads configuration
func (c *Config) LoadConfig() (err error) {
	if err = c.lookupUserConfig(); err != nil {
		return err
	}

	err = c.lookupProjectConfig()

	if err != nil && strings.Contains(err.Error(), "no such file or directory") == false {
		return err
	}

	return nil
}

// FindProjectPathInJSON check if project path in the json
func (c *Config) FindProjectPathInJSON(pc projectSettings) {
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

// GetProjectMainContainer gets project main container
func (c *Config) GetProjectMainContainer() string {
	return c.projectConfig.GetMainContainer()
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
