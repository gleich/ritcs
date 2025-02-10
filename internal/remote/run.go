package remote

import (
	"fmt"
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

	start := time.Now()
	command := strings.Join(cmd, " ")
	fmt.Println()
	timber.Info(fmt.Sprintf("running command \"%s\"", command))
	out, err := session.CombinedOutput(fmt.Sprintf("cd %s && %s", dir, command))
	if err != nil {
		return fmt.Errorf("%v failed to run %s with an output of\n%s", err, cmd, out)
	}
	timber.Done(fmt.Sprintf("finished running in %s", time.Since(start)))
	return nil
}
