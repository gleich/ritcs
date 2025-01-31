package main

import (
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	"pkg.mattglei.ch/timber"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		timber.Fatal(err, "failed to get home directory")
	}

	keyPath := filepath.Join(home, ".ssh", "id_ed25519")
	key, err := os.ReadFile(keyPath)
	if err != nil {
		timber.Fatal(err, "failed to read from", keyPath)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		timber.Fatal(err, "failed to parse private key")
	}

	config := &ssh.ClientConfig{
		User:            "mwg2345",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", "glados.cs.rit.edu:22", config)
	if err != nil {
		timber.Fatal(err, "failed to connect to machine")
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		timber.Fatal(err, "failed to create new session")
	}
	defer session.Close()

	output, err := session.CombinedOutput("ls -la")
	if err != nil {
		timber.Fatal(err, "failed to run remote command")
	}
	timber.Debug(string(output))
}
