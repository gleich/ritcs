package remote

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"pkg.mattglei.ch/ritcs/internal/conf"
	"pkg.mattglei.ch/timber"
)

func RunCmd(client *ssh.Client, config conf.Config, dir string, cmd []string) error {
	session, err := client.NewSession()
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

	start := time.Now()
	command := strings.Join(cmd, " ")
	if !config.Silent {
		fmt.Println()
		timber.Info(fmt.Sprintf("running command \"%s\"", command))
	}
	err = session.Run(fmt.Sprintf("cd %s && %s", dir, command))
	if err != nil {
		return fmt.Errorf("%v failed to run %s", err, cmd)
	}
	if !config.Silent {
		timber.Done(fmt.Sprintf("finished running in %s", time.Since(start)))
	}
	return nil
}
