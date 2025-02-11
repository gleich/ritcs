package cmds

import (
	"fmt"
	"strings"

	"pkg.mattglei.ch/ritcs/internal/conf"
	"pkg.mattglei.ch/ritcs/internal/remote"
	"pkg.mattglei.ch/timber"
)

func Run(cmd []string) error {
	config, err := conf.Load()
	if err != nil {
		timber.Fatal(err, "failed to load configuration file")
	}

	ignoreStatements, err := conf.ReadIgnore(config)
	if err != nil {
		timber.Fatal(err, "failed to read ignore file")
	}

	sshClient, sftpClient, err := remote.EstablishConnection(config)
	if err != nil {
		timber.Fatal(err, "failed to establish connection")
	}
	defer sftpClient.Close()
	defer sshClient.Close()

	tempDir, err := remote.CreateTempDir(config, sftpClient)
	if err != nil {
		timber.Fatal(err, "failed to create temporary directory on server")
	}

	if !config.SkipUpload {
		err = remote.Upload(sftpClient, config, ignoreStatements, tempDir)
		if err != nil {
			timber.Fatal(err, "failed to copy files from host")
		}
	}

	cmdErr := remote.RunCmd(sshClient, config, tempDir, cmd)

	if !config.Silent {
		fmt.Println()
	}
	if !config.SkipDownload {
		err = remote.Download(sftpClient, config, tempDir)
		if err != nil {
			timber.Fatal(err, "failed to copy files from remote")
		}
	}

	err = remote.RemoveTempDir(sftpClient, tempDir)
	if err != nil {
		timber.Fatal(err, "failed to remove temporary directory")
	}

	if cmdErr != nil {
		fmt.Println()
		timber.Fatal(cmdErr, strings.Join(cmd, " "), "excited with a fail exit code")
	}
	return nil
}
