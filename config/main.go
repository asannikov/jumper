package config

// ProjectSettings contains absolut project settings
type ProjectSettings struct {
	Path         string `json:"path"`
	Name         string `json:"name"`
	PhpContainer string
}

// ProjectConfig contains file config
type ProjectConfig struct {
	PhpContainer string `json:"phpContainer"`
}

// FileSystem file management type
type FileSystem struct {
	FileExists     func(string) (bool, error)
	DirExists      func(string) (bool, error)
	SaveConfigFile func(data interface{}, fileName string) error
	ReadConfigFile func(filename string, configuration interface{}) error
}
