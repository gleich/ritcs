package cmds

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"pkg.mattglei.ch/ritcs/internal/conf"
	"pkg.mattglei.ch/ritcs/internal/remote"
	"pkg.mattglei.ch/ritcs/internal/util"
	"pkg.mattglei.ch/timber"
)

func Run(cmd []string) error {
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
	tarFilename := fmt.Sprintf("%s.tar.gz", filepath.Base(remoteTempDir))
	remoteTarPath := filepath.Join(filepath.Dir(remoteTempDir), tarFilename)

	cleanup := sync.WaitGroup{}
	cleanup.Add(1)
	go func() {
		err := remote.CleanupTempDir(sftpClient, remoteTempDir, remoteTarPath)
		if err != nil {
			timber.Error(err, "failed to cleanup temporary directory")
		}
		cleanup.Done()
	}()

	start := time.Now()
	uploadCount, err := remote.UploadCWD(sftpClient, ignoreStatements, remoteTarPath)
	if err != nil {
		timber.Fatal(err, "failed to upload current working directory as a tar file")
	}
	if uploadCount != 0 {
		err = remote.RunTar(sshClient, remoteTempDir, remoteTarPath, true)
		if err != nil {
			timber.Fatal(err, "failed to extract tar file")
		}
		if !conf.Config.Silent {
			timber.Done("uploaded", uploadCount, "files in", util.FormatDuration(time.Since(start)))
		}
	}

	cmdErr := remote.RunCmd(sshClient, remoteTempDir, cmd)

	if !conf.Config.SkipDownload {
		start = time.Now()
		if !conf.Config.Silent {
			fmt.Println()
			timber.Info("downloading files")
		}
		err = remote.RunTar(sshClient, remoteTempDir, remoteTarPath, false)
		if err != nil {
			timber.Fatal(err, "failed to create tar file on remote")
		}
		downloadedFiles, err := remote.DownloadFromTarball(sftpClient, remoteTarPath)
		if err != nil {
			return fmt.Errorf("%v failed to extract remote tar file", err)
		}
		if !conf.Config.Silent {
			timber.Done(
				"downloaded",
				downloadedFiles,
				"files in",
				util.FormatDuration(time.Since(start)),
			)
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
