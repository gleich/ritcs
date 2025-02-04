package cmds

import (
	"pkg.mattglei.ch/ritcs/internal/conf"
	"pkg.mattglei.ch/ritcs/internal/remote"
	"pkg.mattglei.ch/timber"
)

func Get(args []string) error {
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

	err = remote.RunGet(sshClient, tempDir, args)
	if err != nil {
		timber.Fatal(err, "failed to run get command")
	}

	err = remote.CopyFilesFromDir(sftpClient, tempDir)
	if err != nil {
		timber.Fatal(err, "failed to copy over files")
	}

	err = remote.RemoveTempDir(sftpClient, tempDir)
	if err != nil {
		timber.Fatal(err, "failed to remove temporary directory")
	}
	return nil
}
