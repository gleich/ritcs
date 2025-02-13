package remote

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.mattglei.ch/ritcs/internal/conf"
	"go.mattglei.ch/ritcs/internal/util"
	"go.mattglei.ch/timber"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func Exec(sshClient *ssh.Client, dir string, cmd []string) error {
	start := time.Now()
	command := strings.Join(cmd, " ")
	if !conf.Config.Silent {
		fmt.Println()
		timber.Info(fmt.Sprintf("running command \"%s\"", command))
	}

	session, err := sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("%v failed to create new ssh session", err)
	}
	defer session.Close()

	width, height, err := term.GetSize(0)
	if err != nil {
		return fmt.Errorf("%v failed to get terminal size", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14_400,
		ssh.TTY_OP_OSPEED: 14_400,
	}
	err = session.RequestPty("xterm-256color", height, width, modes)
	if err != nil {
		return fmt.Errorf("%v request for pseudo terminal failed", err)
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	err = session.Run(fmt.Sprintf("cd %s && %s", dir, command))
	if err != nil {
		return fmt.Errorf("%v failed to run %s", err, cmd)
	}
	if !conf.Config.Silent {
		timber.Done(fmt.Sprintf("finished running in %s", util.FormatDuration(time.Since(start))))
	}
	return nil
}
