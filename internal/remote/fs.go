package remote

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/sftp"
	"pkg.mattglei.ch/ritcs/internal/conf"
	"pkg.mattglei.ch/timber"
)

func CopyFilesFromDir(client *sftp.Client, dir string) error {
	walker := client.Walk(dir)
	for walker.Step() {
		if err := walker.Err(); err != nil {
			return err
		}

		remotePath := walker.Path()
		relPath, err := filepath.Rel(dir, remotePath)
		if err != nil {
			return fmt.Errorf("%v failed to get relative path", err)
		}
		localPath := filepath.Join("./", relPath)

		if walker.Stat().IsDir() {
			err := os.MkdirAll(localPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("%v failed to create local directory %s", err, localPath)
			}
		} else {
			remoteFile, err := client.Open(remotePath)
			if err != nil {
				return fmt.Errorf("%v failed to open remote file %s", err, localPath)
			}
			defer remoteFile.Close()

			err = os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
			if err != nil {
				return fmt.Errorf("%v failed to create local directory", err)
			}
			localFile, err := os.Create(localPath)
			if err != nil {
				return fmt.Errorf("%v failed to create local file %s", err, localPath)
			}
			defer localFile.Close()

			_, err = io.Copy(localFile, remoteFile)
			if err != nil {
				return fmt.Errorf("%v failed to copy remote file", err)
			}
			timber.Done("copied over", localPath)
		}
	}
	return nil
}

func CreateTempDir(conf conf.Config, client *sftp.Client) (string, error) {
	dir := filepath.Join(conf.Home, "ritcs", strconv.Itoa(rand.Int()))

	err := client.RemoveAll(dir)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("%v failed to remove temporary directory at start", err)
	}

	err = client.MkdirAll(dir)
	if err != nil {
		return "", fmt.Errorf("%v failed to make directory", err)
	}

	return dir, nil
}

func RemoveTempDir(client *sftp.Client, dir string) error {
	err := client.RemoveAll(filepath.Dir(dir))
	if err != nil {
		return fmt.Errorf("%v failed to remove temporary directory %s", err, dir)
	}
	return nil
}
