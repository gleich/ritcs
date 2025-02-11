package cmds

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"pkg.mattglei.ch/ritcs/internal/conf"
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

	remoteTempDir, err := remote.CreateTempDir(sftpClient)
	if err != nil {
		timber.Fatal(err, "failed to create temporary directory on server")
	}
	tarFilename := fmt.Sprintf("%s.tar.gz", strconv.Itoa(rand.Int()))
	remoteTarPath := filepath.Join(filepath.Dir(remoteTempDir), tarFilename)

	cleanup := sync.WaitGroup{}
	cleanup.Add(1)
	go func() {
		err := remote.CleanupTempDir(sftpClient, remoteTempDir)
		if err != nil {
			timber.Error(err, "failed to cleanup temporary directory")
		}
		cleanup.Done()
	}()

	removeUploadedTar := sync.WaitGroup{}
	if !conf.Config.SkipUpload {
		err = remote.UploadCWD(sftpClient, ignoreStatements, remoteTarPath)
		if err != nil {
			timber.Fatal(err, "failed to upload current working directory as a tar file")
		}
		err = remote.RunTar(sshClient, remoteTempDir, remoteTarPath, true)
		if err != nil {
			timber.Fatal(err, "failed to extract tar file")
		}
		removeUploadedTar.Add(1)
		go func() {
			err := sftpClient.Remove(remoteTarPath)
			if err != nil {
				timber.Error(err, "failed to remove tar path", remoteTarPath)
			}
			removeUploadedTar.Done()
		}()
	}

	cmdErr := remote.RunCmd(sshClient, remoteTempDir, cmd)
	removeUploadedTar.Wait()

	if !conf.Config.Silent {
		fmt.Println()
	}
	if !conf.Config.SkipDownload {
		err = remote.RunTar(sshClient, remoteTempDir, remoteTarPath, false)
		if err != nil {
			timber.Fatal(err, "failed to create tar file on remote")
		}
		err = remote.DownloadFromTarball(sftpClient, remoteTarPath)
		if err != nil {
			return fmt.Errorf("%v failed to extract remote tar file", err)
		}
	}

	cleanup.Wait()
	err = sftpClient.Close()
	if err != nil {
		timber.Fatal(err, "failed to close sftp connection")
	}
	err = sshClient.Close()
	if err != nil {
		timber.Fatal(err, "failed to close ssh connection")
	}

	if cmdErr != nil {
		fmt.Println()
		timber.Fatal(cmdErr, strings.Join(cmd, " "), "excited with a fatal exit code")
	}

	return nil
}
