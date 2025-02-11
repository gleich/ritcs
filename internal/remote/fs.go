package remote

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strconv"

	"github.com/pkg/sftp"
	"pkg.mattglei.ch/ritcs/internal/conf"
)

func CreateTempDir(sftpClient *sftp.Client) (string, error) {
	dir := filepath.Join(conf.Config.Home, "ritcs", strconv.Itoa(rand.Int()))
	err := sftpClient.MkdirAll(dir)
	if err != nil {
		return "", fmt.Errorf("%v failed to make directory", err)
	}
	return dir, nil
}

func CleanupTempDir(sftpClient *sftp.Client, tempdir string) error {
	tempDirRoot := filepath.Dir(tempdir)
	fsObjects, err := sftpClient.ReadDir(tempDirRoot)
	if err != nil {
		return fmt.Errorf("%v failed to read %s", err, tempDirRoot)
	}
	for _, obj := range fsObjects {
		path := filepath.Join(tempDirRoot, obj.Name())
		if obj.IsDir() && path != tempdir {
			err = sftpClient.RemoveAll(path)
			if err != nil {
				return fmt.Errorf("%v failed to remove %s", err, path)
			}
		}
	}
	return nil
}

func RemoveTempDir(client *sftp.Client, tempdir string) error {
	err := client.RemoveAll(filepath.Dir(tempdir))
	if err != nil {
		return fmt.Errorf("%v failed to remove temporary directory %s", err, tempdir)
	}
	return nil
}
