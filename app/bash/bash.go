package bash

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// Bash is an object for working with Cli
type Bash struct {
}

// GetCommandLocation get command locaton in docker container
// Use it for detect commad location in related container
func (b Bash) GetCommandLocation() func(string, string) (string, error) {
	return func(containerName string, command string) (string, error) {
		if len(containerName) == 0 {
			return "", errors.New("Container is not defined")
		}

		if len(command) == 0 {
			return "", errors.New("Command is not defined")
		}

		var args = []string{"exec", "-i", containerName, "sh"}

		cmd := exec.Command("docker", args...)

		r, w := io.Pipe()

		cmd.Stdout = w
		cmd.Stderr = w

		var stdin io.WriteCloser
		var err error

		if stdin, err = cmd.StdinPipe(); err != nil {
			return "", fmt.Errorf("error when piping stdin: %w", err)
		}

		defer stdin.Close()

		scanMergedOut := bufio.NewScanner(r)

		fmt.Fprintln(stdin, "which "+command)

		if err = cmd.Start(); err != nil {
			return "", fmt.Errorf("Cannot start command: %w", err)
		}

		if !scanMergedOut.Scan() {
			return "", errors.New("nothing found in Scanner output")
		}

		if scanMergedOut.Err() != nil {
			return "", fmt.Errorf("error while reading scanMergedOut: %s", scanMergedOut.Err())
		}

		return strings.Trim(string(scanMergedOut.Bytes()), " "), nil
	}
}
