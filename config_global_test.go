package main

import (
	"encoding/json"
	"mgt/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluateProjectPath(t *testing.T) {
	curerntPath, _ := os.Getwd()

	DLG := dialogInternal{}
	DLG.addProjectName = func() (string, error) {
		return "project_name", nil
	}
	DLG.addProjectPath = func(path string) (string, error) {
		return curerntPath, nil
	}
	DLG.selectProject = func(p []string) (int, string, error) {
		return -1, "new", nil
	}

	cfg := config.MgtConfig{
		FileSystem: config.FileSystem{
			FileExists: func(string) (bool, error) {
				return true, nil
			},
			DirExists: func(string) (bool, error) {
				return true, nil
			},
			ReadConfigFile: func(filename string, configuration interface{}) error {
				data := `{
					"phpContainer": "phpfpm"
				}`
				return json.Unmarshal([]byte(data), &configuration)
			},
			SaveConfigFile: func(data interface{}, fileName string) error {
				return nil
			},
		},
	}

	getProjectPath := func() (string, error) {
		projectPath, err := cfg.GetProjectPath(func() error {
			return cfg.GenerateConfig(func() error {
				return cfg.EvaluateProjectPath(func(projects []string) (int, string, error) {
					return DLG.selectProject(projects)
				}, func() (string, error) {
					return DLG.addProjectName()
				}, func(path string) (string, error) {
					return DLG.addProjectPath(path)
				})
			})
		})
		return projectPath, err
	}

	projectPath, err := getProjectPath()
	assert.EqualValues(t, curerntPath, projectPath)
	assert.EqualValues(t, nil, err)
}
