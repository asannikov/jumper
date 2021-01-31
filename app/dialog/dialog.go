package dialog

// Dialog contains methods for the iteraction with promptui
type Dialog struct {
	setMainContaner               func([]string) (int, string, error)
	setStartCommand               func() (string, error)
	setStartDocker                func() (string, error)
	setDockerService              func() (string, error)
	setDockerProjectPath          func(string) (string, error)
	setDockerCliXdebugIniFilePath func(string) (string, error)
	setDockerFmpXdebugIniFilePath func(string) (string, error)
	setXdebugFileConfigLocation   func() (int, string, error)
	setDockerShell                func() (int, string, error)

	// Project management
	setSelectProject  func([]string) (int, string, error)
	setAddProjectPath func(string) (string, error)
	setAddProjectName func() (string, error)
}

// InitDialogFunctions initiate all methods
func InitDialogFunctions() Dialog {
	return Dialog{
		setSelectProject:  selectProject,
		setAddProjectPath: addProjectPath,
		setAddProjectName: addProjectName,

		setMainContaner:               setMainContaner,
		setStartCommand:               setStartCommand,
		setStartDocker:                startDocker,
		setDockerService:              dockerService,
		setDockerProjectPath:          dockerProjectPath,
		setDockerShell:                dockerShell,
		setDockerCliXdebugIniFilePath: dockerCliXdebugIniFilePath,
		setDockerFmpXdebugIniFilePath: dockerFmpXdebugIniFilePath,
		setXdebugFileConfigLocation:   xdebugFileConfigLocation,
	}
}

// SetSelectProjectTest is used only for testing
func (d *Dialog) SetSelectProjectTest(f func([]string) (int, string, error)) {
	d.setSelectProject = f
}

// SetAddProjectPathTest is used only for testing
func (d *Dialog) SetAddProjectPathTest(f func(string) (string, error)) {
	d.setAddProjectPath = f
}

// SetAddProjectNameTest is used only for testing
func (d *Dialog) SetAddProjectNameTest(f func() (string, error)) {
	d.setAddProjectName = f
}

// ProjectConfig is for mocking parent functon in main file
type ProjectConfig interface {
	GetProjectName() string
	GetProjectPath() string
	SetProjectName(string)
	SetProjectPath(string)
}
