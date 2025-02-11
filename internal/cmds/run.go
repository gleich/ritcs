package cmds

import (
	"fmt"
	"strings"
	"sync"

	"pkg.mattglei.ch/ritcs/internal/conf"
	"pkg.mattglei.ch/ritcs/internal/host"
	"pkg.mattglei.ch/ritcs/internal/remote"
	"pkg.mattglei.ch/timber"
)

func Run(cmd []string) error {
	err := conf.Load()
	if err != nil {
		timber.Fatal(err, "failed to load configuration file")
	}

	ignoreStatements, err := conf.ReadIgnore()
	if err != nil {
		timber.Fatal(err, "failed to read ignore file")
	}

	sshClient, sftpClient, err := remote.EstablishConnection()
	if err != nil {
		timber.Fatal(err, "failed to establish connection")
	}
	defer sftpClient.Close()
	defer sshClient.Close()

	tempdir, err := remote.CreateTempDir(sftpClient)
	if err != nil {
		timber.Fatal(err, "failed to create temporary directory on server")
	}

	wg := sync.WaitGroup{}
	if !conf.Config.SkipUpload {
		tarpath, err := host.CreateTarball(ignoreStatements)
		if err != nil {
			timber.Fatal(err, "failed to create tarball")
		}
		err = remote.CopyTarball(sftpClient, tempdir, tarpath)
		if err != nil {
			timber.Fatal(err, "failed to copy tarball to remote machine")
		}
		err = remote.RunTar(sshClient, tempdir, tarpath, true)
		if err != nil {
			return fmt.Errorf("%v failed to extract tar file", err)
		}
		go func() {
			err := remote.RemoveTarball(sftpClient, tarpath, &wg)
			if err != nil {
				timber.Error(err, "failed to remove", tarpath, "from remote")
			}
		}()
	}

	cmdErr := remote.RunCmd(sshClient, tempdir, cmd)
	wg.Wait()

	if !conf.Config.Silent {
		fmt.Println()
	}
	if !conf.Config.SkipDownload {
		// err = remote.Download(sftpClient, config, tempDir)
		// if err != nil {
		// 	timber.Fatal(err, "failed to copy files from remote")
		// }
	}

	err = remote.RemoveTempDir(sftpClient, tempdir)
	if err != nil {
		timber.Fatal(err, "failed to remove temporary directory")
	}

	if cmdErr != nil {
		fmt.Println()
		timber.Fatal(cmdErr, strings.Join(cmd, " "), "excited with a fail exit code")
	}
	return nil
}
