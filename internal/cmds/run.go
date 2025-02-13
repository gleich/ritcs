package cmds

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"go.mattglei.ch/ritcs/internal/conf"
	"go.mattglei.ch/ritcs/internal/remote"
	"go.mattglei.ch/ritcs/internal/util"
	"go.mattglei.ch/timber"
	"golang.org/x/crypto/ssh"
)

func Run(cmd []string) {
	err := checkRsyncInstall()
	if err != nil {
		timber.Fatal(err, "failed to check for rsync installation")
	}

	cwd, err := os.Getwd()
	if err != nil {
		timber.Fatal(err, "failed to get working directory")
	}
	projectPath := remote.ProjectPath(cwd)

	var (
		sshClient  *ssh.Client
		sshSession *ssh.Session
	)
	sshConnection := sync.WaitGroup{}
	sshConnection.Add(1)
	go func() {
		sshClient, err = remote.EstablishConnection()
		if err != nil {
			timber.Fatal(err, "failed to connect to remote host machine")
		}
		sshSession, err = remote.CreateSession(sshClient)
		if err != nil {
			timber.Fatal(err, "failed to create ssh session")
		}
		sshConnection.Done()
	}()

	s := spinner.New(util.CustomSpinner, 10*time.Millisecond)
	s.Suffix = fmt.Sprintf(" uploading files to %s", conf.Config.Host)
	if !conf.Config.Silent {
		s.Start()
	}
	err = remote.RunRsync(projectPath, remote.Upload)
	if err != nil {
		timber.Fatal(err, "failed to run rsync")
	}
	if !conf.Config.Silent {
		s.Stop()
		timber.Done("uploaded files to", conf.Config.Host)
	}

	sshConnection.Wait()
	err = remote.Exec(sshSession, projectPath, cmd)
	if err != nil {
		timber.Fatal(err, "failed to execute command")
	}

	closeSSH := sync.WaitGroup{}
	closeSSH.Add(1)
	go func() {
		err := sshClient.Close()
		if err != nil {
			timber.Fatal(err, "failed to close ssh connection")
		}
		closeSSH.Done()
	}()

	if !conf.Config.SkipDownload {
		if !conf.Config.Silent {
			fmt.Println()
			s.Suffix = fmt.Sprintf(" downloading files from %s", conf.Config.Host)
			s.Start()
		}
		err = remote.RunRsync(projectPath, remote.Download)
		if err != nil {
			timber.Fatal(err, "failed to download files using rsync")
		}
		if !conf.Config.Silent {
			s.Stop()
			timber.Done("downloaded files from", conf.Config.Host)
		}
	}

	closeSSH.Wait()
}

func checkRsyncInstall() error {
	_, err := exec.LookPath("rsync")
	if errors.Is(err, exec.ErrNotFound) {
		timber.FatalMsg("rsync binary not found. please install rsync as it is required by ritcs")
	}
	if err != nil {
		return fmt.Errorf("%v failed to check to see if rsync is installed", err)
	}
	return nil
}
