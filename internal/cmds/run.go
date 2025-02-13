package cmds

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"go.mattglei.ch/ritcs/internal/conf"
	"go.mattglei.ch/ritcs/internal/remote"
	"go.mattglei.ch/timber"
)

func Run(cmd []string) error {
	err := checkRsyncInstall()
	if err != nil {
		return fmt.Errorf("%v failed to check to see if rsync is installed", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("%v failed to get working directory", err)
	}
	projectPath := remote.ProjectPath(cwd)

	sshClient, err := remote.EstablishConnection()
	if err != nil {
		return fmt.Errorf("%v failed to establish connection to remote machine", err)
	}
	defer sshClient.Close()

	frames := []string{
		"[    ]", // 0%
		"[=   ]", // 6.25%
		"[==  ]", // 12.5%
		"[=== ]", // 18.75%
		"[ ===]", // 25%
		"[  ==]", // 31.25%
		"[   =]", // 37.5%
		"[    ]", // 43.75%
		"[   =]", // 50%
		"[  ==]", // 56.25%
		"[ ===]", // 62.5%
		"[====]", // 68.75%
		"[=== ]", // 75%
		"[==  ]", // 81.25%
		"[=   ]", // 87.5%
		"[    ]", // 93.75% / 100%
	}
	s := spinner.New(frames, 20*time.Millisecond)
	s.Suffix = fmt.Sprintf(" uploading files to %s", conf.Config.Host)
	if !conf.Config.Silent {
		s.Start()
	}
	err = remote.RunRsync(projectPath, remote.Upload)
	if err != nil {
		return fmt.Errorf("%v failed to run rsync", err)
	}
	if !conf.Config.Silent {
		s.Stop()
		timber.Done("uploaded files to", conf.Config.Host)
	}

	err = remote.Exec(sshClient, projectPath, cmd)
	if err != nil {
		return fmt.Errorf("%v failed to run command", err)
	}

	if !conf.Config.SkipDownload {
		if !conf.Config.Silent {
			fmt.Println()
			s.Suffix = fmt.Sprintf(" downloading files from %s", conf.Config.Host)
			s.Start()
		}
		err = remote.RunRsync(projectPath, remote.Download)
		if err != nil {
			return fmt.Errorf("%v failed to download files using rsync", err)
		}
		if !conf.Config.Silent {
			s.Stop()
			timber.Done("downloaded files from", conf.Config.Host)
		}
	}

	return nil
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
