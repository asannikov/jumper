package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func fileExist(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, err
	}
	return !info.IsDir(), nil
}

func dirExists(path string) (bool, error) {
	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			return true, nil
		}
		return false, fmt.Errorf("Path %s is a file ", path)
	} else if os.IsNotExist(err) {
		// path does *not* exist
		return false, err
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false, err
	}
}

func readConfigFile(filename string, configuration interface{}) (err error) {
	if len(filename) == 0 {
		return fmt.Errorf("Reading config error, file name is empty")
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), &configuration)
}

func saveConfigFile(data interface{}, fileName string) error {
	file, _ := json.MarshalIndent(data, "", " ")
	err := ioutil.WriteFile(fileName, file, 0644)
	return err
}
