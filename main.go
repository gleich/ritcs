package main

import (
	"os"
	"time"

	"pkg.mattglei.ch/timber"
)

func main() {
	timber.SetTimezone(time.Local)
	timber.SetTimeFormat("03:04:05")

	if len(os.Args) <= 1 {
		timber.FatalMsg("please provide arguments to get command")
	}

	if os.Args[1] == "setup" {
		err := setup()
		if err != nil {
			timber.Fatal(err, "failed to setup user")
		}
		return
	}

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

	tempDir, err := createTempDir(conf, sftpClient)
	if err != nil {
		timber.Fatal(err, "failed to create temporary directory on server")
	}

	err = runGet(sshClient, tempDir)
	if err != nil {
		timber.Fatal(err, "failed to run get command")
	}

	err = copyFilesOver(sftpClient, tempDir)
	if err != nil {
		timber.Fatal(err, "failed to copy over files")
	}

	err = removeTempDir(sftpClient, tempDir)
	if err != nil {
		timber.Fatal(err, "Failed to remove temporary directory")
	}
}
