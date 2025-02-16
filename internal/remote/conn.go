package remote

import (
	"fmt"
	"os"
	"path/filepath"

	"go.mattglei.ch/ritcs/internal/conf"
	"golang.org/x/crypto/ssh"
)

func EstablishConnection() (*ssh.Client, error) {
	key, err := os.ReadFile(conf.Config.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("%w failed to read from key path %s", err, conf.Config.KeyPath)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("%w failed to parse private key", err)
	}

	config := &ssh.ClientConfig{
		User:            filepath.Base(conf.Config.Home),
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", conf.Config.Host, conf.Config.Port),
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("%w failed to create connection", err)
	}

	return sshClient, nil
}
