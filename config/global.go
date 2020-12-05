package config

// GlobalConfig contains file config
type GlobalConfig struct {
	Projects []ProjectSettings `json:"projects"`
}

// AddNewProject adds new project
func (g *GlobalConfig) AddNewProject(p ProjectSettings) {
	g.Projects = append(g.Projects, p)
}

// GetProjectNameList gets project name list
func (g *GlobalConfig) GetProjectNameList() []string {
	pl := []string{}

	for _, p := range g.Projects {
		if p.Path != "" {
			pl = append(pl, p.Name)
		}
	}

	return pl
}

// FindProjectPathInJSON find project path in json
func (g *GlobalConfig) FindProjectPathInJSON(f func(ProjectSettings) (bool, error)) error {
	for _, p := range g.Projects {
		if r, err := f(p); err == nil && r == true {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}
