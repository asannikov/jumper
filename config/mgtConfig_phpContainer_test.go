package config

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPhpContainer(t *testing.T) {
	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return false, fmt.Errorf("Cannot open file")
			},
		},
	}

	// case 1
	containerName, err := cfg.GetPhpContainer(func() error {
		return fmt.Errorf("Generate config error")
	})

	assert.Equal(t, "", containerName)
	assert.EqualError(t, err, "Generate config error")

	// case 2
	cfg.FileSystem.ReadConfigFile = func(filename string, configuration interface{}) error {
		return fmt.Errorf("Cannot read config")
	}

	cfg.FileSystem.FileExists = func(string) (bool, error) {
		return true, nil
	}

	containerName, err = cfg.GetPhpContainer(func() error {
		return nil
	})

	assert.Equal(t, "", containerName)
	assert.EqualError(t, err, "Cannot read config")

	// case 3
	cfg.FileSystem.ReadConfigFile = func(filename string, configuration interface{}) error {
		return nil
	}

	cfg.FileSystem.FileExists = func(string) (bool, error) {
		return true, nil
	}

	containerName, err = cfg.GetPhpContainer(func() error {
		return nil
	})

	assert.Equal(t, "", containerName)
	assert.EqualError(t, err, "Php container missing in config file")

	// case 4
	cfg.FileSystem.ReadConfigFile = func(filename string, configuration interface{}) error {
		data := `{
	"phpContainer": "phpfpm"
}`
		return json.Unmarshal([]byte(data), &configuration)
	}

	containerName, err = cfg.GetPhpContainer(func() error {
		return nil
	})

	assert.Equal(t, "phpfpm", containerName)
	assert.Equal(t, nil, err)
}
func TestEvaluatePhpContainer(t *testing.T) {

	cfg := MgtConfig{
		FileSystem: FileSystem{
			FileExists: func(string) (bool, error) {
				return true, nil
			},
			DirExists: func(string) (bool, error) {
				return true, nil
			},
			ReadConfigFile: func(filename string, configuration interface{}) error {
				return nil
			},
			SaveConfigFile: func(data interface{}, fileName string) error {
				return nil
			},
		},
	}

	// case 1
	err := cfg.EvaluatePhpContainer(func() (int, string, error) {
		return 1, "dockerPhpContainer", nil
	})

	assert.Equal(t, nil, err)

	// case 2
	cfg.FileSystem.FileExists = func(string) (bool, error) {
		return false, nil
	}

	err = cfg.EvaluatePhpContainer(func() (int, string, error) {
		return 1, "dockerPhpContainer", nil
	})

	assert.Equal(t, nil, err)

	// case 3
	cfg.FileSystem.FileExists = func(string) (bool, error) {
		return false, fmt.Errorf("Cannot open file")
	}

	err = cfg.EvaluatePhpContainer(func() (int, string, error) {
		return 1, "dockerPhpContainer", nil
	})

	assert.EqualError(t, err, "Cannot open file")

	// case 4
	cfg.FileSystem.FileExists = func(string) (bool, error) {
		return true, nil
	}

	cfg.FileSystem.ReadConfigFile = func(filename string, configuration interface{}) error {
		return fmt.Errorf("Cannot read file")
	}

	err = cfg.EvaluatePhpContainer(func() (int, string, error) {
		return 1, "dockerPhpContainer", nil
	})

	assert.EqualError(t, err, "Cannot read file")

	// case 5
	cfg.FileSystem.FileExists = func(string) (bool, error) {
		return false, nil
	}

	err = cfg.EvaluatePhpContainer(func() (int, string, error) {
		return 1, "dockerPhpContainer", nil
	})

	assert.Equal(t, nil, err)

	// case 5
	cfg.FileSystem.FileExists = func(string) (bool, error) {
		return true, nil
	}

	cfg.FileSystem.ReadConfigFile = func(filename string, configuration interface{}) error {
		return nil
	}

	err = cfg.EvaluatePhpContainer(func() (int, string, error) {
		return 0, "", fmt.Errorf("promtui error")
	})

	assert.EqualError(t, err, "promtui error")

	// case 6
	cfg.FileSystem.FileExists = func(string) (bool, error) {
		return true, nil
	}

	cfg.FileSystem.ReadConfigFile = func(filename string, configuration interface{}) error {
		return nil
	}

	cfg.FileSystem.SaveConfigFile = func(data interface{}, fileName string) error {
		return fmt.Errorf("cannot save file")
	}

	err = cfg.EvaluatePhpContainer(func() (int, string, error) {
		return 0, "dockerPhpContainer", nil
	})

	assert.EqualError(t, err, "cannot save file")

}
