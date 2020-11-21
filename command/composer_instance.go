package command

import (
	"bufio"
	"bytes"
	"fmt" // imports as package "cli"
	"io"
	"log"
	"os/exec"
)

// Composer is an instance of related containers
type Composer struct {
	phpInstance      string
	composerInstance string
}

type cliInstance struct {
	stdin         io.WriteCloser
	stdMergedOut  io.ReadCloser
	scanMergedOut *bufio.Scanner
}

// CallComposerCommand1 generates Composer command (docker-compose exec phpfpm composer)
/**
 * @todo do composer install/update with no memory limit
 * https://stackoverflow.com/questions/43099116/error-the-input-device-is-not-a-tty
 */

// GetContainerNames gets php and composer container names
func (c *Composer) GetContainerNames(containerPhp string) (string, string, error) {
	var phpInstance, composerInstance string

	var err error
	if phpInstance, composerInstance, err = c.getInstance(containerPhp); err != nil {
		return "", "", err
	}

	return phpInstance, composerInstance, nil
}

func (c *Composer) getInstance(containerPhp string) (string, string, error) {
	var binary = "docker"
	var args = []string{"exec", "-i", containerPhp, "bash"}

	p := cliInstance{}

	cmd := exec.Command(binary, args...)

	r, w := io.Pipe()

	p.stdMergedOut = r

	cmd.Stdout = w
	cmd.Stderr = w

	var err error

	if p.stdin, err = cmd.StdinPipe(); err != nil {
		return "", "", fmt.Errorf("error when piping stdin: %w", err)
	}

	p.scanMergedOut = bufio.NewScanner(r)

	//p.scanMergedOut.Split(splitReadyToken)

	err = cmd.Start()

	if err != nil {
		return "", "", fmt.Errorf("Cannot start command: %w", err)
	}

	var phpInstance, composerInstance string

	if phpInstance, err = p.getPhpInstance(); err != nil {
		return "", "", fmt.Errorf("no php instance found: %w", err)
	}

	if composerInstance, err = p.getComposerInstance(); err != nil {
		return "", "", fmt.Errorf("no composer instance found: %w", err)
	}

	p.stdMergedOut.Close()
	p.stdin.Close()

	return phpInstance, composerInstance, nil
}

func (c *cliInstance) getPhpInstance() (string, error) {
	fmt.Fprintln(c.stdin, "which php")

	if !c.scanMergedOut.Scan() {
		return "", fmt.Errorf("nothing on stdMergedOut")
	}

	if c.scanMergedOut.Err() != nil {
		return "", fmt.Errorf("error while reading stdMergedOut: %s", c.scanMergedOut.Err())
	}

	return string(c.scanMergedOut.Bytes()), nil
}

func (c *cliInstance) getComposerInstance() (string, error) {
	fmt.Fprintln(c.stdin, "which composer")

	if !c.scanMergedOut.Scan() {
		return "", fmt.Errorf("nothing on stdMergedOut")
	}

	if c.scanMergedOut.Err() != nil {
		return "", fmt.Errorf("error while reading stdMergedOut: %s", c.scanMergedOut.Err())
	}

	return string(c.scanMergedOut.Bytes()), nil
}

func splitReadyToken(data []byte, atEOF bool) (int, []byte, error) {

	var readyToken = []byte("$\r\n")
	var readyTokenLen = len(readyToken)

	log.Println(string(data))
	idx := bytes.Index(data, readyToken)
	log.Println(atEOF, len(data))
	if idx == -1 {
		if atEOF && len(data) > 0 {
			return 0, data, fmt.Errorf("no final token found")
		}

		return 0, nil, nil
	}

	return idx + readyTokenLen, data[:idx], nil
}
