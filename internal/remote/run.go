package remote

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"pkg.mattglei.ch/timber"
)

func RunCmd(client *ssh.Client, dir string, cmd []string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("%v failed to create new ssh session", err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14_400,
		ssh.TTY_OP_OSPEED: 14_400,
	}
	err = session.RequestPty("xterm-256color", 80, 40, modes)
	if err != nil {
		return fmt.Errorf("%v request for pseudo terminal failed", err)
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	start := time.Now()
	command := strings.Join(cmd, " ")
	fmt.Println()
	timber.Info(fmt.Sprintf("running command \"%s\"", command))
	err = session.Run(fmt.Sprintf("cd %s && %s", dir, command))
	if err != nil {
		return fmt.Errorf("%v failed to run %s", err, cmd)
	}
	timber.Done(fmt.Sprintf("finished running in %s", time.Since(start)))
	return nil
}
