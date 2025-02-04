package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Home    string `toml:"home,required"`
	Host    string `toml:"host,required"`
	Port    int    `toml:"port"`
	KeyPath string `toml:"key_path,required"`
}

func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("%v failed to get user's home directory", err)
	}
	return filepath.Join(home, ".config", "ritcs", "config.toml"), nil
}

func Load() (Config, error) {
	path, err := Path()
	if err != nil {
		return Config{}, fmt.Errorf("%v failed to get configuration path", err)
	}

	var conf Config
	_, err = toml.DecodeFile(path, &conf)
	if err != nil {
		return Config{}, fmt.Errorf("%v failed to decode TOML config file from %s", err, path)
	}
	return conf, nil
}
