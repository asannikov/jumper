package config

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProjectPath(t *testing.T) {
	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return true, nil
			},
			ReadConfigFile: func(filename string, configuration interface{}) error {
				data := `{
	"phpContainer": "phpfpm"
}`
				return json.Unmarshal([]byte(data), &configuration)
			},
		},
	}

	// case 1
	path, err := cfg.GetProjectPath(func() error {
		return nil
	})

	assert.Equal(t, "", path)
	assert.Equal(t, nil, err)

	// case 2
	path, err = cfg.GetProjectPath(func() error {
		return fmt.Errorf("Cannot generate config")
	})
	assert.Equal(t, "", path)
	assert.EqualError(t, err, "Cannot generate config")

	// case 3

	cfg.FileSystem.FileExists = func(string) (bool, error) {
		return false, fmt.Errorf("Cannot open file")
	}

	path, err = cfg.GetProjectPath(func() error {
		return nil
	})

	assert.Equal(t, "", path)
	assert.EqualError(t, err, "Cannot open file")
}

func TestProjectListManagementDialogCase1(t *testing.T) {
	cpd := createProjectDialog{}

	err := projectListManagementDialog(func() (string, error) {
		return "project_name", nil
	}, func(path string) (string, error) {
		return "project_path", nil
	}, &cpd)

	assert.Equal(t, "project_name", cpd.projectName)
	assert.Equal(t, "project_path", cpd.projectPath)
	assert.Equal(t, nil, err)
}

func TestProjectListManagementDialogCase2(t *testing.T) {
	cpd := createProjectDialog{}

	err := projectListManagementDialog(func() (string, error) {
		return "project_name", fmt.Errorf("Name cannot be empty")
	}, func(path string) (string, error) {
		return "project_path", nil
	}, &cpd)

	assert.Equal(t, "", cpd.projectName)
	assert.Equal(t, "", cpd.projectPath)
	assert.EqualError(t, err, "Name cannot be empty")
}

func TestProjectListManagementDialogCase3(t *testing.T) {
	cpd := createProjectDialog{}

	err := projectListManagementDialog(func() (string, error) {
		return "project_name", nil
	}, func(path string) (string, error) {
		return "project_path", fmt.Errorf("Path cannot be empty")
	}, &cpd)

	assert.Equal(t, "project_name", cpd.projectName)
	assert.Equal(t, "", cpd.projectPath)
	assert.EqualError(t, err, "Path cannot be empty")
}

func TestProjectListManagementDialogCase4(t *testing.T) {
	cpd := createProjectDialog{
		projectName: "old_name",
	}

	err := projectListManagementDialog(func() (string, error) {
		return "new_name", nil
	}, func(path string) (string, error) {
		return "project_path", nil
	}, &cpd)

	assert.Equal(t, "old_name", cpd.projectName)
	assert.Equal(t, "project_path", cpd.projectPath)
	assert.Equal(t, nil, err)
}

func TestEvaluateProjectPathCase1(t *testing.T) {
	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return false, fmt.Errorf("Cannot open file")
			},
		},
	}

	err := cfg.EvaluateProjectPath(func(projects []string) (int, string, error) {
		return 0, "project_name", nil
	}, func() (string, error) {
		return "new_name", nil
	}, func(path string) (string, error) {
		return "project_path", nil
	})

	assert.EqualError(t, err, "Cannot open file")
}

func TestEvaluateProjectPathCase2(t *testing.T) {
	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return true, fmt.Errorf("Cannot open file")
			},
		},
	}

	err := cfg.EvaluateProjectPath(func(projects []string) (int, string, error) {
		return 0, "project_name", nil
	}, func() (string, error) {
		return "new_name", nil
	}, func(path string) (string, error) {
		return "project_path", nil
	})

	assert.EqualError(t, err, "Cannot open file")
}

func TestEvaluateProjectPathCase3(t *testing.T) {
	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return true, nil
			},
			ReadConfigFile: func(filename string, configuration interface{}) error {
				return fmt.Errorf("Cannot read config")
			},
		},
	}

	err := cfg.EvaluateProjectPath(func(projects []string) (int, string, error) {
		return 0, "project_name", nil
	}, func() (string, error) {
		return "new_name", nil
	}, func(path string) (string, error) {
		return "project_path", nil
	})

	assert.EqualError(t, err, "Cannot read config")
}

func TestEvaluateProjectPathCase4(t *testing.T) {
	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return true, nil
			},
			ReadConfigFile: func(filename string, configuration interface{}) error {
				return nil
			},
		},
	}

	err := cfg.EvaluateProjectPath(func(projects []string) (int, string, error) {
		return 0, "project_name", nil
	}, func() (string, error) {
		return "", fmt.Errorf("Promtui error")
	}, func(path string) (string, error) {
		return "", nil
	})

	assert.EqualError(t, err, "Promtui error")
}

func TestEvaluateProjectPathCase5(t *testing.T) {
	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return true, nil
			},
			ReadConfigFile: func(filename string, configuration interface{}) error {
				data := `{
	"projects": [
		{
			"path": "/path/to/project1",
			"name": "Name0"
		},
		{
			"path": "/path/to/project2",
			"name": "Name1"
		}
	]
}`
				return json.Unmarshal([]byte(data), &configuration)
			},
		},
	}

	err := cfg.EvaluateProjectPath(func(projects []string) (int, string, error) {
		return 0, "", fmt.Errorf("Promtui error")
	}, func() (string, error) {
		return "project_name", nil
	}, func(path string) (string, error) {
		return "project_path", nil
	})

	assert.EqualError(t, err, "Promtui error")
}

func TestEvaluateProjectPathCase6(t *testing.T) {
	//currentDir, _ := os.Getwd()

	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return true, nil
			},
			ReadConfigFile: func(filename string, configuration interface{}) error {
				data := `{
	"projects": [
		{
			"path": "/path/to/project1",
			"name": "Name0"
		},
		{
			"path": "/path/to/project2",
			"name": "Name1"
		}
	]
}`
				return json.Unmarshal([]byte(data), &configuration)
			},
		},
	}

	err := cfg.EvaluateProjectPath(func(projects []string) (int, string, error) {
		return -1, "new_project_name", nil
	}, func() (string, error) {
		return "", nil
	}, func(path string) (string, error) {
		return "", fmt.Errorf("Promtui error")
	})

	assert.EqualError(t, err, "Promtui error")
}

func TestEvaluateProjectPathCase7(t *testing.T) {
	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return true, nil
			},
			ReadConfigFile: func(filename string, configuration interface{}) error {
				data := `{
	"projects": [
		{
			"path": "/path/to/project1",
			"name": "Name0"
		},
		{
			"path": "/path/to/project2",
			"name": "project_name"
		}
	]
}`
				return json.Unmarshal([]byte(data), &configuration)
			},
			DirExists: func(string) (bool, error) {
				return false, fmt.Errorf("Directory does not exist")
			},
		},
	}

	err := cfg.EvaluateProjectPath(func(projects []string) (int, string, error) {
		return 0, "project_name", nil
	}, func() (string, error) {
		return "", fmt.Errorf("Promtui error1")
	}, func(path string) (string, error) {
		return "", fmt.Errorf("Promtui error2")
	})

	assert.EqualError(t, err, "Directory does not exist")
}

func TestEvaluateProjectPathCase8(t *testing.T) {
	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return true, nil
			},
			ReadConfigFile: func(filename string, configuration interface{}) error {
				data := `{
	"projects": [
		{
			"path": "/path/to/project1",
			"name": "Name0"
		},
		{
			"path": "/path/to/project2",
			"name": "project_name"
		}
	]
}`
				return json.Unmarshal([]byte(data), &configuration)
			},
			DirExists: func(string) (bool, error) {
				return true, fmt.Errorf("Directory does not exist")
			},
		},
	}

	err := cfg.EvaluateProjectPath(func(projects []string) (int, string, error) {
		return 0, "project_name", nil
	}, func() (string, error) {
		return "", fmt.Errorf("Promtui error1")
	}, func(path string) (string, error) {
		return "", fmt.Errorf("Promtui error2")
	})

	assert.EqualError(t, err, "Directory does not exist")
}

func TestEvaluateProjectPathCase9(t *testing.T) {
	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return true, nil
			},
			ReadConfigFile: func(filename string, configuration interface{}) error {
				data := `{
	"projects": [
		{
			"path": "/path/to/project1",
			"name": "Name0"
		},
		{
			"path": "",
			"name": "project_name"
		}
	]
}`
				return json.Unmarshal([]byte(data), &configuration)
			},
			DirExists: func(string) (bool, error) {
				return true, nil
			},
		},
	}

	err := cfg.EvaluateProjectPath(func(projects []string) (int, string, error) {
		return 0, "project_name", nil
	}, func() (string, error) {
		return "", fmt.Errorf("Promtui error1")
	}, func(path string) (string, error) {
		return "", fmt.Errorf("Promtui error2")
	})

	assert.EqualError(t, err, "Cannot use project in the root path")
}

func TestEvaluateProjectPathCase10(t *testing.T) {
	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return true, nil
			},
			ReadConfigFile: func(filename string, configuration interface{}) error {
				data := `{
	"projects": [
		{
			"path": "/path/to/project1",
			"name": "Name0"
		},
		{
			"path": "/path/to/project2/",
			"name": "project_name"
		}
	]
}`
				return json.Unmarshal([]byte(data), &configuration)
			},
			DirExists: func(string) (bool, error) {
				return true, nil
			},
		},
	}

	err := cfg.EvaluateProjectPath(func(projects []string) (int, string, error) {
		return 0, "project_name", nil
	}, func() (string, error) {
		return "", fmt.Errorf("Promtui error1")
	}, func(path string) (string, error) {
		return "", fmt.Errorf("Promtui error2")
	})

	assert.EqualError(t, err, "chdir /path/to/project2: no such file or directory")
}

func TestEvaluateProjectPathCase11(t *testing.T) {
	currentDir, _ := os.Getwd()

	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return true, nil
			},
			ReadConfigFile: func(filename string, configuration interface{}) error {
				data := `{
	"projects": [
		{
			"path": "/path/to/project1",
			"name": "Name0"
		},
		{
			"path": "` + currentDir + `/",
			"name": "project_name"
		}
	]
}`
				return json.Unmarshal([]byte(data), &configuration)
			},
			DirExists: func(string) (bool, error) {
				return true, nil
			},
			SaveConfigFile: func(data interface{}, fileName string) error {
				return fmt.Errorf("Cannot save config")
			},
		},
	}

	err := cfg.EvaluateProjectPath(func(projects []string) (int, string, error) {
		return 0, "project_name", nil
	}, func() (string, error) {
		return "", fmt.Errorf("Promtui error1")
	}, func(path string) (string, error) {
		return "", fmt.Errorf("Promtui error2")
	})

	assert.EqualError(t, err, "Cannot save config")
}

func TestEvaluateProjectPathCase12(t *testing.T) {
	currentDir, _ := os.Getwd()

	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return true, nil
			},
			ReadConfigFile: func(filename string, configuration interface{}) error {
				data := `{
	"projects": [
		{
			"path": "/path/to/project1",
			"name": "Name0"
		},
		{
			"path": "` + currentDir + `/",
			"name": "project_name"
		}
	]
}`
				return json.Unmarshal([]byte(data), &configuration)
			},
			DirExists: func(string) (bool, error) {
				return true, nil
			},
			SaveConfigFile: func(data interface{}, fileName string) error {
				return nil
			},
		},
	}

	err := cfg.EvaluateProjectPath(func(projects []string) (int, string, error) {
		return 0, "project_name", nil
	}, func() (string, error) {
		return "", fmt.Errorf("Promtui error1")
	}, func(path string) (string, error) {
		return "", fmt.Errorf("Promtui error2")
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "project_name", cfg.projectSettings.Name)
	assert.Equal(t, currentDir+"/", cfg.projectSettings.Path)
}

func TestEvaluateProjectPathCase13(t *testing.T) {
	currentDir, _ := os.Getwd()

	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return true, nil
			},
			ReadConfigFile: func(filename string, configuration interface{}) error {
				data := `{
	"projects": [
		{
			"path": "/path/to/project1",
			"name": "Name0"
		},
		{
			"path": "/path/to/project2/",
			"name": "Name1"
		}
	]
}`
				return json.Unmarshal([]byte(data), &configuration)
			},
			DirExists: func(string) (bool, error) {
				return true, nil
			},
			SaveConfigFile: func(data interface{}, fileName string) error {
				return nil
			},
		},
	}

	err := cfg.EvaluateProjectPath(func(projects []string) (int, string, error) {
		return -1, "new_project_name", nil
	}, func() (string, error) {
		return "new_project_name", nil
	}, func(path string) (string, error) {
		return currentDir, nil
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(cfg.globalSettings.Projects))
}
