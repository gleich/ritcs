package main

import (
	"errors"
	"fmt"
	"io/fs"
	"math/rand/v2"
	"path/filepath"
	"strconv"

	"github.com/pkg/sftp"
	"pkg.mattglei.ch/timber"
)

func createTempDir(conf config, client *sftp.Client) (string, error) {
	dir := filepath.Join(*conf.Home, "ritcsget", strconv.Itoa(rand.Int()))

	err := client.RemoveAll(dir)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("%v failed to remove temporary directory at start", err)
	}

	err = client.MkdirAll(dir)
	if err != nil {
		return "", fmt.Errorf("%v failed to make directory", err)
	}

	timber.Done("crated temporary directory in", dir)
	return dir, nil
}
