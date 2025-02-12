package cmds

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/huh"
	"go.mattglei.ch/ritcs/internal/conf"
	"go.mattglei.ch/timber"
)

func Setup() error {
	config := conf.Configuration{Port: 22}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("home").
				Description("Home directory for your CS user account.").
				Value(&config.Home),
			huh.NewInput().
				Title("host").
				Description("Hostname of the CS machine you want to ssh into").
				Placeholder("glados.cs.rit.edu").
				Value(&config.Host),
			huh.NewInput().
				Title("key_path").
				Description("Local path to the private key used to authenticate with the CS machine.").
				Placeholder(".../id_rsa").
				Value(&config.KeyPath).
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

	b, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("%v failed to marshal config into toml", err)
	}

	path, err := conf.Path()
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
