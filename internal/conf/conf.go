package conf

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"go.mattglei.ch/timber"
)

var Config Configuration

type Configuration struct {
	Home         string `toml:"home,required"`
	Host         string `toml:"host,required"`
	KeyPath      string `toml:"key_path,required"`
	Port         int    `toml:"port"`
	SkipDownload bool   `toml:"skip_download"`
	Silent       bool   `toml:"silent"`
	User         string `toml:"user,omitempty"`
}

func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("%w failed to get user's home directory", err)
	}
	return filepath.Join(home, ".config", "ritcs", "config.toml"), nil
}

func Load() error {
	path, err := Path()
	if err != nil {
		return fmt.Errorf("%w failed to get configuration path", err)
	}

	_, err = os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		timber.FatalMsg("no configuration file found. please run 'ritcs setup'")
	}

	var conf Configuration
	_, err = toml.DecodeFile(path, &conf)
	if err != nil {
		return fmt.Errorf(
			"%v failed to decode TOML config file from %s",
			err,
			path,
		)
	}

	conf.User = filepath.Base(conf.Home)

	Config = conf
	return nil
}
