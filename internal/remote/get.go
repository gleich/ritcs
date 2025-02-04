package remote

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

func RunGet(client *ssh.Client, dir string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("%v failed to create new ssh session", err)
	}
	defer session.Close()

	cmd := fmt.Sprintf("get %s", strings.Join(os.Args[1:], " "))
	out, err := session.CombinedOutput(fmt.Sprintf("cd %s && %s", dir, cmd))
	if err != nil {
		return fmt.Errorf("%v failed to run %s with an output of\n%s", err, cmd, out)
	}

	return nil
}
