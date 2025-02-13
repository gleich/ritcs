package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
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
		return "", fmt.Errorf("%v failed to get user's home directory", err)
	}
	return filepath.Join(home, ".config", "ritcs", "config.toml"), nil
}

func Load() error {
	path, err := Path()
	if err != nil {
		return fmt.Errorf("%v failed to get configuration path", err)
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
