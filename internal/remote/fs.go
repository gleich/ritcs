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

func CopyFilesFromHost(client *sftp.Client, tempDir string) error {
	var files, folders int

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("%v failed to get working directory", err)
	}

	outputNewline := false
	err = filepath.Walk(cwd, func(localPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skipping symlinks
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		relPath, err := filepath.Rel(cwd, localPath)
		if err != nil {
			return fmt.Errorf("%v failed to get relative path for %s", err, localPath)
		}
		remotePath := filepath.Join(tempDir, relPath)

		if info.IsDir() {
			if err := client.MkdirAll(remotePath); err != nil {
				return fmt.Errorf("%v failed to create remote directory %s", err, remotePath)
			}
			folders++
		} else {
			localFile, err := os.Open(localPath)
			if err != nil {
				return fmt.Errorf("%v failed to open local file %s", err, localPath)
			}

			remoteDir := filepath.Dir(remotePath)
			if err := client.MkdirAll(remoteDir); err != nil {
				return fmt.Errorf("%v failed to create remote directory %s", err, remoteDir)
			}

			remoteFile, err := client.Create(remotePath)
			if err != nil {
				return fmt.Errorf("%v failed to create remote file %s", err, remotePath)
			}

			if _, err := io.Copy(remoteFile, localFile); err != nil {
				return fmt.Errorf("%v failed to copy local file %s to remote", err, localPath)
			}

			err = localFile.Close()
			if err != nil {
				return fmt.Errorf("%v failed to close local file", err)
			}
			err = remoteFile.Close()
			if err != nil {
				return fmt.Errorf("%v failed to close remote file", err)
			}

			if !outputNewline {
				fmt.Println()
				outputNewline = true
			}
			timber.Done("uploaded", relPath)
			files++
		}

		perm := info.Mode().Perm()
		err = client.Chmod(remotePath, perm)
		if err != nil {
			return fmt.Errorf("%v failed to set mode of %s for %s", err, perm, remotePath)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("%v failed to walk local directory", err)
	}
	return nil
}

func CopyFilesFromRemote(client *sftp.Client, tempDir string) error {
	walker := client.Walk(tempDir)
	var files, folders int
	for walker.Step() {
		if err := walker.Err(); err != nil {
			return err
		}

		remotePath := walker.Path()
		relPath, err := filepath.Rel(tempDir, remotePath)
		if err != nil {
			return fmt.Errorf("%v failed to get relative path", err)
		}
		localPath := filepath.Join("./", relPath)

		if walker.Stat().IsDir() {
			err := os.MkdirAll(localPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("%v failed to create local directory %s", err, localPath)
			}
			folders++
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

			timber.Done("downloaded", relPath)
			files++
		}

		perm := walker.Stat().Mode().Perm()
		err = os.Chmod(localPath, perm)
		if err != nil {
			return fmt.Errorf("%v failed to set %s for %s", err, perm, localPath)
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
