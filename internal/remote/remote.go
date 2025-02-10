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

func EstablishConnection(conf conf.Config) (*ssh.Client, *sftp.Client, error) {
	key, err := os.ReadFile(conf.KeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("%v failed to read from key path %s", err, conf.KeyPath)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, nil, fmt.Errorf("%v failed to parse private key", err)
	}

	config := &ssh.ClientConfig{
		User:            filepath.Base(conf.Home),
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port), config)
	if err != nil {
		return nil, nil, fmt.Errorf("%v failed to create connection", err)
	}
	timber.Done("established SSH connection to", conf.Host)

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return nil, nil, fmt.Errorf("%v failed to create sftp client", err)
	}
	timber.Done("established SFTP connection to", conf.Host)

	return conn, sftpClient, nil
}
