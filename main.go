package main

import (
	"os"

	"pkg.mattglei.ch/timber"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		timber.Fatal(err, "failed to get user's home directory")
	}

	conf, err := loadConfig(home)
	if err != nil {
		timber.Fatal(err, "failed to load configuration file")
	}

	client, session, err := establishConnection(home, conf)
	if err != nil {
		timber.Fatal(err, "failed to establish connection")
	}
	defer client.Close()
	defer session.Close()

	output, err := session.CombinedOutput("ls -la")
	if err != nil {
		timber.Fatal(err, "failed to run remote command")
	}
	timber.Debug(string(output))
}
