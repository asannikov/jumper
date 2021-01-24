package docker

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
)

type cliApp interface {
	GetReader() io.Reader
	GetWriter() io.Writer
	GetErrWriter() io.Writer
}

// Exec run command using docker api
func (d *Docker) Exec(container string, cnf *types.ExecConfig, ca cliApp) (status int, err error) {
	status = 1

	fmt.Printf("\ncommand: %s\n\n", strings.Join(cnf.Cmd, " "))

	ctx := context.Background()
	r, err := d.GetClient().ContainerExecCreate(ctx, container, *cnf)

	if err != nil {
		return 1, err
	}

	stream, err := d.GetClient().ContainerExecAttach(ctx, r.ID, types.ExecStartCheck{})

	if err != nil {
		return 1, err
	}

	outputErr := make(chan error)

	go func() {
		var err error
		if cnf.Tty {
			_, err = io.Copy(ca.GetWriter(), stream.Reader)
		} else {
			_, err = stdcopy.StdCopy(ca.GetWriter(), ca.GetErrWriter(), stream.Reader)
		}
		outputErr <- err
	}()

	go func() {
		defer stream.CloseWrite()
		io.Copy(ca.GetWriter(), stream.Conn)
	}()

	for {
		inspect, err := d.GetClient().ContainerExecInspect(ctx, r.ID)
		if err != nil {
			return 0, err
		}
		if !inspect.Running {
			status = inspect.ExitCode
			break
		}
		time.Sleep(time.Second)
	}

	err = <-outputErr
	return
}
