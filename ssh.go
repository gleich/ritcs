package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func establishConnection(home string, conf config) (*ssh.Client, *ssh.Session, error) {
	key, err := os.ReadFile(*conf.KeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("%v failed to read from key path %s", err, *conf.KeyPath)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, nil, fmt.Errorf("%v failed to parse private key", err)
	}

	config := &ssh.ClientConfig{
		User:            *conf.Username,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", *conf.Host, conf.Port), config)
	if err != nil {
		return nil, nil, fmt.Errorf("%v failed to create connection", err)
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("%v failed to create new session", err)
	}

	return conn, session, nil
}
