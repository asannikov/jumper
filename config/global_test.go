package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindProjectPathInJSON(t *testing.T) {
	gc := GlobalConfig{
		Projects: []ProjectSettings{
			ProjectSettings{
				Path: "path",
				Name: "name",
			},
		},
	}

	var foundPath string

	projectName := "name"
	err := gc.FindProjectPathInJSON(func(p ProjectSettings) (bool, error) {
		if p.Name == projectName {
			foundPath = p.Path
			return true, nil
		}
		return false, nil
	})

	assert.EqualValues(t, "path", foundPath)
	assert.EqualValues(t, nil, err)
}

func TestFindProjectPathInJSONEmpty(t *testing.T) {
	gc := GlobalConfig{
		Projects: []ProjectSettings{},
	}

	var foundPath string

	projectName := "name"
	err := gc.FindProjectPathInJSON(func(p ProjectSettings) (bool, error) {
		if p.Name == projectName {
			foundPath = p.Path
			return true, nil
		}
		return false, nil
	})

	assert.EqualValues(t, "", foundPath)
	assert.EqualValues(t, nil, err)
}

func TestFindProjectPathInJSONError(t *testing.T) {
	gc := GlobalConfig{
		Projects: []ProjectSettings{
			ProjectSettings{
				Path: "path",
				Name: "name",
			},
		},
	}

	err := gc.FindProjectPathInJSON(func(p ProjectSettings) (bool, error) {
		return false, fmt.Errorf("There is an error")
	})

	assert.EqualError(t, err, "There is an error")

	err = gc.FindProjectPathInJSON(func(p ProjectSettings) (bool, error) {
		return true, fmt.Errorf("There is an error")
	})

	assert.EqualError(t, err, "There is an error")
}
