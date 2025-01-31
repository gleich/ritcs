package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/huh"
	"pkg.mattglei.ch/timber"
)

type config struct {
	Home    string `toml:"home,required"`
	Host    string `toml:"host,required"`
	Port    int    `toml:"port"`
	KeyPath string `toml:"key_path,required"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("%v failed to get user's home directory", err)
	}
	return filepath.Join(home, ".config", "ritcsget", "config.toml"), nil
}

func loadConfig() (config, error) {
	path, err := configPath()
	if err != nil {
		return config{}, fmt.Errorf("%v failed to get configuration path", err)
	}

	var conf config
	_, err = toml.DecodeFile(path, &conf)
	if err != nil {
		return config{}, fmt.Errorf("%v failed to decode TOML config file from %s", err, path)
	}
	return conf, nil
}

func setup() error {
	conf := config{Port: 22}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("home").
				Description("Home directory for your CS user account.").
				Value(&conf.Home),
			huh.NewInput().
				Title("host").
				Description("Hostname of the CS machine you want to ssh into").
				Placeholder("glados.cs.rit.edu").
				Value(&conf.Host),
			huh.NewInput().
				Title("key_path").
				Description("Local path to the private key used to authenticate with the CS machine.").
				Placeholder(".../id_rsa").
				Value(&conf.KeyPath).
				Validate(func(s string) error {
					_, err := os.Stat(s)
					if err != nil {
						return err
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeBase())

	err := form.Run()
	if err != nil {
		return fmt.Errorf("%v failed to ask user for configuration", err)
	}

	b, err := toml.Marshal(conf)
	if err != nil {
		return fmt.Errorf("%v failed to marshal config into toml", err)
	}

	path, err := configPath()
	if err != nil {
		return fmt.Errorf("%v failed to get configuration path", err)
	}

	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return fmt.Errorf("%v failed to create folder", err)
	}

	err = os.WriteFile(path, b, 0650)
	if err != nil {
		return fmt.Errorf("%v failed to write configuration to file", err)
	}

	timber.Done("created configuration at", path)
	return nil
}
