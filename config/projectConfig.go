package config

// ProjectConfig contains project settings
type ProjectConfig struct {
	Path              string `json:"-"`
	Name              string `json:"name"`
	MainContainer     string `json:"main_container"`
	StartCommand      string `json:"start_command"`
	DockerProjectPath string `json:"path"`
}

// GetPath gets path to project
func (p *ProjectConfig) GetPath() string {
	return p.Path
}

// GetDockerProjectPath gets path to project
func (p *ProjectConfig) GetDockerProjectPath() string {
	return p.DockerProjectPath
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
