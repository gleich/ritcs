package remote

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"pkg.mattglei.ch/ritcs/internal/conf"
	"pkg.mattglei.ch/timber"
)

func EstablishConnection() (*ssh.Client, *sftp.Client, error) {
	key, err := os.ReadFile(conf.Config.KeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("%v failed to read from key path %s", err, conf.Config.KeyPath)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, nil, fmt.Errorf("%v failed to parse private key", err)
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
		return nil, nil, fmt.Errorf("%v failed to create connection", err)
	}
	if !conf.Config.Silent {
		timber.Done("established SSH connection to", conf.Config.Host)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, nil, fmt.Errorf("%v failed to create sftp client", err)
	}
	if !conf.Config.Silent {
		timber.Done("established SFTP connection to", conf.Config.Host)
	}

	return sshClient, sftpClient, nil
}
