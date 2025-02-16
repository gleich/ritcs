package remote

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"go.mattglei.ch/ritcs/internal/conf"
	"go.mattglei.ch/ritcs/internal/util"
	"go.mattglei.ch/timber"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func Exec(session *ssh.Session, dir string, cmd []string) error {
	defer session.Close()
	start := time.Now()
	command := strings.Join(cmd, " ")
	if !conf.Config.Silent {
		fmt.Println()
		timber.Info(fmt.Sprintf("running command \"%s\"", command))
	}
	err := session.Run(fmt.Sprintf("cd %s && %s", dir, command))
	if err != nil {
		return fmt.Errorf("%w failed to run %s", err, cmd)
	}
	if !conf.Config.Silent {
		timber.Done(fmt.Sprintf("finished running in %s", util.FormatDuration(time.Since(start))))
	}
	return nil
}

func CreateSession(sshClient *ssh.Client) (*ssh.Session, error) {
	session, err := sshClient.NewSession()
	if err != nil {
		return nil, fmt.Errorf("%w failed to create new ssh session", err)
	}

	width, height, err := term.GetSize(0)
	if err != nil {
		return nil, fmt.Errorf("%w failed to get terminal size", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14_400,
		ssh.TTY_OP_OSPEED: 14_400,
	}
	err = session.RequestPty("xterm-256color", height, width, modes)
	if err != nil {
		return nil, fmt.Errorf("%w request for pseudo terminal failed", err)
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	return session, nil
}

type rsyncOperation int

const (
	Download rsyncOperation = iota
	Upload   rsyncOperation = iota
)

func RunRsync(projectPath string, op rsyncOperation) error {
	var (
		args           []string
		remoteLocation = fmt.Sprintf("%s@%s:%s/", conf.Config.User, conf.Config.Host, projectPath)
	)

	switch op {
	case Upload:
		args = []string{
			"-ahrzW",
			"--perms",
			"--rsync-path", fmt.Sprintf("mkdir -p %s && rsync", projectPath), // create parent directory
			"--omit-dir-times",
			"--exclude=.git",
			"--exclude=.DS_Store",
			"./",
			remoteLocation,
			"--delete",
		}
	case Download:
		args = []string{
			"-ahrzW",
			"--omit-dir-times",
			"--perms",
			remoteLocation,
			".",
		}
	}

	cmd := exec.Command("rsync", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w failed to run rsync\n%s", err, string(out))
	}

	return nil
}
