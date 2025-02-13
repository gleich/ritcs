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
	"go.mattglei.ch/ritcs/internal/util"
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

	sshClient, sftpClient, err := remote.EstablishConnection()
	if err != nil {
		return fmt.Errorf("%v failed to establish connection to remote machine", err)
	}
	defer sftpClient.Close()
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
	start := time.Now()
	s := spinner.New(frames, 20*time.Millisecond)
	s.Suffix = fmt.Sprintf(" uploading files to %s", conf.Config.Host)
	fmt.Println()
	s.Start()
	err = sftpClient.MkdirAll(projectPath)
	if err != nil {
		return fmt.Errorf("%v failed to create project path %s", err, projectPath)
	}
	err = remote.RunRsync(projectPath, remote.Upload)
	if err != nil {
		return fmt.Errorf("%v failed to run rsync", err)
	}
	s.Stop()
	timber.Done("uploaded files to", conf.Config.Host, "in", util.FormatDuration(time.Since(start)))

	err = remote.Exec(sshClient, projectPath, cmd)
	if err != nil {
		return fmt.Errorf("%v failed to run command", err)
	}

	if !conf.Config.SkipDownload {
		fmt.Println()
		start = time.Now()
		s.Suffix = fmt.Sprintf(" downloading files from %s", conf.Config.Host)
		s.Start()
		err = remote.RunRsync(projectPath, remote.Download)
		if err != nil {
			return fmt.Errorf("%v failed to download files using rsync", err)
		}
		s.Stop()
		timber.Done(
			"downloaded files from",
			conf.Config.Host,
			"in",
			util.FormatDuration(time.Since(start)),
		)

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
