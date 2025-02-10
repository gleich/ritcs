package cmds

import (
	"fmt"

	"pkg.mattglei.ch/ritcs/internal/conf"
	"pkg.mattglei.ch/ritcs/internal/remote"
	"pkg.mattglei.ch/timber"
)

func Run(cmd []string) error {
	conf, err := conf.Load()
	if err != nil {
		timber.Fatal(err, "failed to load configuration file")
	}

	sshClient, sftpClient, err := remote.EstablishConnection(conf)
	if err != nil {
		timber.Fatal(err, "failed to establish connection")
	}
	defer sftpClient.Close()
	defer sshClient.Close()

	tempDir, err := remote.CreateTempDir(conf, sftpClient)
	if err != nil {
		timber.Fatal(err, "failed to create temporary directory on server")
	}

	err = remote.CopyFilesFromHost(sftpClient, tempDir)
	if err != nil {
		timber.Fatal(err, "failed to copy files from host")
	}

	err = remote.RunCmd(sshClient, tempDir, cmd)
	if err != nil {
		timber.Fatal(err, "failed to run get command")
	}

	fmt.Println()
	err = remote.CopyFilesFromRemote(sftpClient, tempDir)
	if err != nil {
		timber.Fatal(err, "failed to copy files from remote")
	}

	err = remote.RemoveTempDir(sftpClient, tempDir)
	if err != nil {
		timber.Fatal(err, "failed to remove temporary directory")
	}
	return nil
}
