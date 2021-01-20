package docker

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
)

// Exec run command using docker api
func (d *Docker) Exec(container string, commands []string) error {

	//defer d.GetClient().Close()
	ctx := context.Background()
	r, err := d.GetClient().ContainerExecCreate(ctx, container, types.ExecConfig{
		Tty:          true,
		AttachStderr: true,
		AttachStdout: true,
		AttachStdin:  true,
		Detach:       true,
		Cmd:          commands,
	})

	if err != nil {
		return err
	}

	out, err := d.GetClient().ContainerLogs(ctx, container, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}

	attach, err := d.GetClient().ContainerExecAttach(ctx, r.ID, types.ExecStartCheck{})

	if err != nil {
		return err
	}

	defer attach.Close()

	err = d.GetClient().ContainerExecStart(ctx, r.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	go io.Copy(os.Stdout, out)

	<-make(chan struct{})
	//log.Println(attach.Reader.)
	return err
}
