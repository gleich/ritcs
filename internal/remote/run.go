package remote

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"pkg.mattglei.ch/ritcs/internal/conf"
	"pkg.mattglei.ch/ritcs/internal/util"
	"pkg.mattglei.ch/timber"
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

func ExecTar(sshClient *ssh.Client, remoteTempDir, remoteTarPath string, extract bool) error {
	session, err := sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("%v failed to create new ssh session", err)
	}
	defer session.Close()

	var cmd string
	if extract {
		cmd = fmt.Sprintf("tar -xzvf %s -C %s", remoteTarPath, remoteTempDir)
	} else {
		cmd = fmt.Sprintf("tar -czvf %s --warning=no-file-changed .", remoteTarPath)
	}

	out, err := session.CombinedOutput(fmt.Sprintf("cd %s && %s", remoteTempDir, cmd))
	if err != nil {
		return fmt.Errorf("%v failed to run %s: %s", err, cmd, string(out))
	}
	return nil
}
