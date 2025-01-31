package main

import (
	"time"

	"pkg.mattglei.ch/timber"
)

func main() {
	timber.SetTimezone(time.Local)
	timber.SetTimeFormat("03:04:05")

	conf, err := loadConfig()
	if err != nil {
		timber.Fatal(err, "failed to load configuration file")
	}

	sshClient, sftpClient, err := establishConnection(conf)
	if err != nil {
		timber.Fatal(err, "failed to establish connection")
	}
	defer sftpClient.Close()
	defer sshClient.Close()

	_, err = createTempDir(conf, sftpClient)
	if err != nil {
		timber.Fatal(err, "failed to create temporary directory on server")
	}
}
