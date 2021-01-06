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

// projectSettings is not the same as ProjectConfig, but the similar one.
// it helps to call FindProjectPathInJSON from the outside (ie main function), where
// projectSettings is used as projectConfig, see main_start.go
type projectSettings interface {
	GetProjectName() string
	GetProjectPath() string
	SetProjectName(string)
	SetProjectPath(string)
}

// SetProjectPath set project path for config, it's not the same as projectSettings
func (c *Config) SetProjectPath(path string) {
	c.projectConfig.Path = path
}

// SetProjectName set project path for config, it's not the same as projectSettings
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
	if err = c.fileSystem.ReadConfigFile(c.UserFile, c.globalConfig); err == nil {
		return nil
	}

	if err != nil && strings.Contains(err.Error(), "no such file or directory") == true {
		err = c.fileSystem.SaveConfigFile(c.globalConfig, c.UserFile)
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

// GetShell gets shell command
func (c *Config) GetShell() string {
	return c.projectConfig.GetShell()
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

	return c.fileSystem.SaveConfigFile(c.globalConfig, c.UserFile)
}

// main container

// SaveShellCommand saves linux shell command
func (c *Config) SaveShellCommand(cmd string) (err error) {
	c.projectConfig.Shell = cmd
	return c.saveProjectFile()
}

// SaveContainerNameToProjectConfig saves container name into project file
func (c *Config) SaveContainerNameToProjectConfig(cn string) (err error) {
	c.projectConfig.MainContainer = cn
	return c.saveProjectFile()
}

// GetProjectMainContainer gets project main container
func (c *Config) GetProjectMainContainer() string {
	return c.projectConfig.GetMainContainer()
}

// Docker project path

// GetProjectDockerPath gets project main container
func (c *Config) GetProjectDockerPath() string {
	return c.projectConfig.GetDockerProjectPath()
}

// SaveDockerProjectPath saves path to project in container into project file
func (c *Config) SaveDockerProjectPath(path string) (err error) {
	c.projectConfig.DockerProjectPath = path
	return c.saveProjectFile()
}

// Xdebug

// GetXDebugConfigLocaton gets cli xdebug ini file path
func (c *Config) GetXDebugConfigLocaton() string {
	return c.projectConfig.GetXDebugConfigLocaton()
}

// GetXDebugCliIniPath gets cli xdebug ini file path
func (c *Config) GetXDebugCliIniPath() string {
	return c.projectConfig.GetXDebugCliIniPath()
}

// GetXDebugFpmIniPath gets fpm xdebug ini file path
func (c *Config) GetXDebugFpmIniPath() string {
	return c.projectConfig.GetXDebugFpmIniPath()
}

// SaveXDebugConifgLocaton saves xdebug file location
func (c *Config) SaveXDebugConifgLocaton(path string) (err error) {
	c.projectConfig.XDebugLocation = path
	return c.saveProjectFile()
}

// SaveDockerCliXdebugIniFilePath saves xdebug cli ini file path into project file
func (c *Config) SaveDockerCliXdebugIniFilePath(path string) (err error) {
	c.projectConfig.XDebugCliIniPath = path
	return c.saveProjectFile()
}

// SaveDockerFpmXdebugIniFilePath saves xdebug fpm ini file path into project file
func (c *Config) SaveDockerFpmXdebugIniFilePath(path string) (err error) {
	c.projectConfig.XDebugFpmIniPath = path
	return c.saveProjectFile()
}

// main start command

// GetStartCommand gets start command
func (c *Config) GetStartCommand() string {
	return c.projectConfig.GetStartCommand()
}

// SaveStartCommandToProjectConfig saves container name into project file
func (c *Config) SaveStartCommandToProjectConfig(cmd string) (err error) {
	c.projectConfig.StartCommand = cmd
	return c.saveProjectFile()
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

// Copyright

// EnableCopyright Enable copyright output
func (c *Config) EnableCopyright() error {
	c.globalConfig.EnableCopyright()
	return c.fileSystem.SaveConfigFile(c.globalConfig, c.UserFile)
}

// DisableCopyright Disable copyright output
func (c *Config) DisableCopyright() error {
	c.globalConfig.DisableCopyright()
	return c.fileSystem.SaveConfigFile(c.globalConfig, c.UserFile)
}

// ShowCopyrightText check the status of copyright output
func (c *Config) ShowCopyrightText() bool {
	return c.globalConfig.ShowCopyrightText()
}

// Docker instance command

// SetDockerCommand define docker command
func (c *Config) SetDockerCommand(command string) error {
	c.globalConfig.SetDockerCommand(command)
	return c.fileSystem.SaveConfigFile(c.globalConfig, c.UserFile)
}

// GetDockerCommand gets the docker command
func (c *Config) GetDockerCommand() string {
	return c.globalConfig.GetDockerCommand()
}
