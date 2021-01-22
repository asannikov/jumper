package docker

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
)

// Exec run command using docker api
func (d *Docker) Exec(container string, commands []string) (status int, err error) {

	status = 1

	fmt.Println(strings.Join(commands, " "))

	//defer d.GetClient().Close()
	ctx := context.Background()
	r, err := d.GetClient().ContainerExecCreate(ctx, container, types.ExecConfig{
		User:         "limesoda",
		Tty:          false,
		AttachStderr: true,
		AttachStdout: true,
		AttachStdin:  true,
		Detach:       true,
		Cmd:          commands,
		WorkingDir:   "/var/www",
	})

	if err != nil {
		return 0, err
	}

	stream, err := d.GetClient().ContainerExecAttach(ctx, r.ID, types.ExecStartCheck{})

	if err != nil {
		return 0, err
	}

	outputErr := make(chan error)

	go func() {
		var err error
		if false { // tty
			_, err = io.Copy(os.Stdout, stream.Reader)
		} else {
			_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, stream.Reader)
		}
		outputErr <- err
	}()

	go func() {
		defer stream.CloseWrite()
		io.Copy(stream.Conn, os.Stdout)
	}()

	/*if cfg.Tty {
		_, winCh, _ := sess.Pty()
		go func() {
			for win := range winCh {
				err := docker.ContainerExecResize(ctx, eresp.ID, types.ResizeOptions{
					Height: uint(win.Height),
					Width:  uint(win.Width),
				})
				if err != nil {
					log.Println(err)
					break
				}
			}
		}()
	}*/
	for {
		inspect, err := d.GetClient().ContainerExecInspect(ctx, r.ID)
		if err != nil {
			log.Println(err)
		}
		if !inspect.Running {
			status = inspect.ExitCode
			break
		}
		time.Sleep(time.Second)
	}

	return
	//err = <-outputErr
	//data, err := ioutil.ReadAll(stream.Reader)
	//fmt.Println(string(data))

	/*
		err = d.GetClient().ContainerExecStart(ctx, r.ID, types.ExecStartCheck{
			Detach: false,
			Tty:    false,
		})

		if err != nil {
			return err
		}

		go io.Copy(os.Stdout, attach.Reader)

		<-make(chan struct{})
		//log.Println(attach.Reader.)
		return err */
}
