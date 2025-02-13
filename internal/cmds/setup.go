package cmds

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/huh"
	"go.mattglei.ch/ritcs/internal/conf"
	"go.mattglei.ch/ritcs/internal/remote"
	"go.mattglei.ch/timber"
)

func Setup() {
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
		timber.Fatal(err, "failed to ask user for configuration")
	}

	conf.Config = config
	sshClient, err := remote.EstablishConnection()
	if err != nil {
		timber.FatalMsg("connection test failed:", err.Error())
	}
	defer sshClient.Close()
	timber.Done("connection test PASSED")

	b, err := toml.Marshal(config)
	if err != nil {
		timber.Fatal(err, "failed to marshal config into toml")
	}

	path, err := conf.Path()
	if err != nil {
		timber.Fatal(err, "failed to get configuration path")
	}

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		timber.Fatal(err, "failed to create folder")
	}

	err = os.WriteFile(path, b, 0644)
	if err != nil {
		timber.Fatal(err, "failed to write configuration to file")
	}

	timber.Done("created configuration at", path)
}
