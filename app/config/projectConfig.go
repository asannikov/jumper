package config

/**
 * @todo add different sh types bash/sh
 */

// ProjectConfig contains project settings
type ProjectConfig struct {
	Path              string `json:"-"`
	Name              string `json:"name"`
	MainContainer     string `json:"main_container"`
	StartCommand      string `json:"start_command"`
	DockerProjectPath string `json:"path"`
	XDebugLocation    string `json:"xdebug_location"`
	XDebugCliIniPath  string `json:"xdebug_path_cli"`
	XDebugFpmIniPath  string `json:"xdebug_path_fpm"`
	Shell             string `json:"shell"`
	ContainerUser     string `json:"main_container_user"`
}

// GetShell gets path to project
func (p *ProjectConfig) GetShell() string {
	return p.Shell
}

// GetPath gets path to project
func (p *ProjectConfig) GetPath() string {
	return p.Path
}

// GetDockerProjectPath gets path to project
func (p *ProjectConfig) GetDockerProjectPath() string {
	return p.DockerProjectPath
}

// GetXDebugCliIniPath gets path to project
func (p *ProjectConfig) GetXDebugCliIniPath() string {
	return p.XDebugCliIniPath
}

// GetXDebugFpmIniPath gets path to project
func (p *ProjectConfig) GetXDebugFpmIniPath() string {
	return p.XDebugFpmIniPath
}

// GetXDebugConfigLocaton gets xdebug file config location
func (p *ProjectConfig) GetXDebugConfigLocaton() string {
	return p.XDebugLocation
}

// GetName gets project name
func (p *ProjectConfig) GetName() string {
	return p.Name
}

// GetMainContainer gets php container name
func (p *ProjectConfig) GetMainContainer() string {
	return p.MainContainer
}

// GetStartCommand gets php container name
func (p *ProjectConfig) GetStartCommand() string {
	return p.StartCommand
}

// GetMainContainerUser gets php container main user
func (p *ProjectConfig) GetMainContainerUser() string {
	return p.ContainerUser
}
