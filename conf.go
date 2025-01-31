package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type config struct {
	Home     *string `toml:"home"`
	Username *string `toml:"username"`
	Host     *string `toml:"host"`
	Port     int     `toml:"port"`
	KeyPath  *string `toml:"key_path"`
}

func loadConfig() (config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return config{}, fmt.Errorf("%v failed to get user's home directory", err)
	}
	path := filepath.Join(home, ".config", "ritcsget", "config.toml")

	var conf config
	_, err = toml.DecodeFile(path, &conf)
	if err != nil {
		return config{}, fmt.Errorf("%v failed to decode TOML config file from %s", err, path)
	}
	if conf.Port == 0 {
		conf.Port = 22
	}
	return conf, nil
}
