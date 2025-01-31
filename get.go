package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"pkg.mattglei.ch/timber"
)

func runGet(client *ssh.Client, tempDir string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("%v failed to create new ssh session", err)
	}
	defer session.Close()

	cmd := fmt.Sprintf("get %s", strings.Join(os.Args[1:], " "))
	timber.Info("running", cmd)
	out, err := session.CombinedOutput(fmt.Sprintf("cd %s && %s", tempDir, cmd))
	if err != nil {
		return fmt.Errorf("%v failed to run %s with an output of\n%s", err, cmd, out)
	}
	timber.Done("finished running", cmd)

	return nil
}
