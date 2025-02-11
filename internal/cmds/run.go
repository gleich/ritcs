package cmds

import (
	"fmt"
	"os"
	"strings"

	"pkg.mattglei.ch/ritcs/internal/conf"
	"pkg.mattglei.ch/ritcs/internal/local"
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
		zipPath, err := local.CreateTarball(config, ignoreStatements)
		if err != nil {
			timber.Fatal(err, "failed to create tarball")
		}
		timber.Debug(zipPath)
		os.Exit(0)
	}

	cmdErr := remote.RunCmd(sshClient, config, tempDir, cmd)

	if !config.Silent {
		fmt.Println()
	}
	if !config.SkipDownload {
		// err = remote.Download(sftpClient, config, tempDir)
		// if err != nil {
		// 	timber.Fatal(err, "failed to copy files from remote")
		// }
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
