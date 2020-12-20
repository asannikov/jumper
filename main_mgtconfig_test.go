package main

/*func getJumperConfig() config.JumperConfig {
	return config.JumperConfig{
		FileName: "test.json",
		FileUser: "/path/to/.test.json",
		FileSystem: config.FileSystem{
			FileExists: func(filename string) (bool, error) {
				return true, nil
			},
			DirExists: func(filename string) (bool, error) {
				return true, nil
			},
			ReadConfigFile: func(filename string, configuration interface{}) error {
				data := `{}`
				return json.Unmarshal([]byte(data), &configuration)
			},
			SaveConfigFile: func(data interface{}, fileName string) error {
				return nil
			},
			GoToProjectPath: func(filepath string) error {
				return nil
			},
		},
	}
}*/

/**
 * - Project File is not found,
 * - error "no such file or directory"
 * - a) handleDialog returns error
 * - b) handleDialog returns nil
 */
/* func TestLoadProjectConfigCase1(t *testing.T) {

	cfg := getJumperConfig()

	c := 0
	cfg.FileSystem.FileExists = func(filename string) (bool, error) {
		if filename == "test.json" && c < 3 {
			c++
			return false, fmt.Errorf("Error: no such file or directory")
		}

		return true, nil
	}

	err := cfg.LoadProjectConfig(func() error {
		return fmt.Errorf("dialog error")
	})

	assert.EqualError(t, err, "dialog error")

	err = cfg.LoadProjectConfig(func() error {
		return nil
	})

	assert.Equal(t, nil, err)
}*/

/**
 * - Project File is not found,
 * - error "no such file or directory"
 * - handleDialog returns error
 */
/*func TestLoadProjectConfigCase2(t *testing.T) {

	cfg := getJumperConfig()

	cfg.FileSystem.FileExists = func(filename string) (bool, error) {
		if filename == "test.json" {
			return false, fmt.Errorf("Error: no such file or directory")
		}

		return true, nil
	}

	err := cfg.LoadProjectConfig(func() error {
		return fmt.Errorf("dialog error")
	})

	assert.EqualError(t, err, "dialog error")
}*/
