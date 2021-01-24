package docker

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
)

// CopyTo syncs data into docker container
func (d *Docker) CopyTo(container string, sourcePath string, dstPath string) (err error) {
	copyToContainerOptions := types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}

	var tar io.Reader

	if tar, err = createTarArchiveFromPath(sourcePath); err != nil {
		return err
	}

	return d.GetClient().CopyToContainer(context.Background(), container, dstPath, tar, copyToContainerOptions)
}

func createTarArchiveFromPath(path string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	ok := filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(strings.Replace(file, path, "", -1), string(filepath.Separator))
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		_, err = io.Copy(tw, f)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
		return nil
	})

	if ok != nil {
		return nil, ok
	}
	ok = tw.Close()
	if ok != nil {
		return nil, ok
	}
	return bufio.NewReader(&buf), nil
}
