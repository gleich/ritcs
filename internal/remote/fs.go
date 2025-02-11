package remote

import (
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"path/filepath"
	"strconv"

	"github.com/pkg/sftp"
	"pkg.mattglei.ch/ritcs/internal/conf"
)

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
