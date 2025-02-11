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
)

func CopyTarball(client *sftp.Client, tempdir string, tarpath string) error {
	remoteLocation := filepath.Join(tempdir, filepath.Base(tarpath))

	localFile, err := os.Open(tarpath)
	if err != nil {
		return fmt.Errorf("%v failed to read from tar file at %s", err, tarpath)
	}
	defer localFile.Close()

	remoteFile, err := client.Create(remoteLocation)
	if err != nil {
		return fmt.Errorf("%v failed to create remote file at %s", err, remoteLocation)
	}
	defer remoteFile.Close()

	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("%v failed to copy local tar file to remote tar file", err)
	}

	return nil
}

func CreateTempDir(client *sftp.Client) (string, error) {
	dir := filepath.Join(conf.Config.Home, "ritcs", strconv.Itoa(rand.Int()))

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
