package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testGlobalFindProjectPathInJSONCase1(t *testing.T) {
	gcfg := GlobalConfig{
		Projects: []GlobalProjectConfig{
			GlobalProjectConfig{
				Path: "path1",
				Name: "name1",
			},
			GlobalProjectConfig{
				Path: "path2",
				Name: "name2",
			},
			GlobalProjectConfig{
				Path: "path3",
				Name: "name3",
			},
		},
	}

	iterator := func(project GlobalProjectConfig) (bool, error) {
		if project.Path == "path2" {
			return false, errors.New("Unknown error")
		}
		return false, nil
	}

	err := gcfg.FindProjectPathInJSON(iterator)
	assert.EqualError(t, err, "Unknown error")
}

func testGlobalFindProjectPathInJSONCase2(t *testing.T) {
	gcfg := GlobalConfig{
		Projects: []GlobalProjectConfig{
			GlobalProjectConfig{
				Path: "path1",
				Name: "name1",
			},
			GlobalProjectConfig{
				Path: "path2",
				Name: "name2",
			},
			GlobalProjectConfig{
				Path: "path3",
				Name: "name3",
			},
		},
	}

	iterator := func(project GlobalProjectConfig) (bool, error) {
		if project.Path == "path2" {
			return true, errors.New("Unknown error")
		}
		return false, nil
	}

	err := gcfg.FindProjectPathInJSON(iterator)
	assert.EqualError(t, err, "Unknown error")
}

func testGlobalFindProjectPathInJSONCase3(t *testing.T) {
	gcfg := GlobalConfig{
		Projects: []GlobalProjectConfig{
			GlobalProjectConfig{
				Path: "path1",
				Name: "name1",
			},
			GlobalProjectConfig{
				Path: "path2",
				Name: "name2",
			},
			GlobalProjectConfig{
				Path: "path3",
				Name: "name3",
			},
		},
	}

	iterator := func(project GlobalProjectConfig) (bool, error) {
		return false, nil
	}

	assert.Equal(t, nil, gcfg.FindProjectPathInJSON(iterator))
}

func testGetProjectNameListCase1(t *testing.T) {
	gcfg := GlobalConfig{
		Projects: []GlobalProjectConfig{
			GlobalProjectConfig{
				Path: "path1",
				Name: "name1",
			},
			GlobalProjectConfig{
				Path: "",
				Name: "name2",
			},
			GlobalProjectConfig{
				Path: "path2",
				Name: "",
			},
			GlobalProjectConfig{
				Path: "path3",
				Name: "name3",
			},
		},
	}

	pl := gcfg.GetProjectNameList()

	assert.Equal(t, 3, len(pl))
}
