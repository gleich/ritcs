package remote

import (
	"fmt"
	"os/exec"

	"go.mattglei.ch/ritcs/internal/conf"
)

type operation int

const (
	Download operation = iota
	Upload   operation = iota
)

func RunRsync(projectPath string, op operation) error {
	var (
		args           []string
		remoteLocation = fmt.Sprintf("%s@%s:%s/", conf.Config.User, conf.Config.Host, projectPath)
	)

	switch op {
	case Upload:
		args = []string{
			"-ahr",
			"--perms",
			"--exclude=.git",
			"--exclude=.DS_Store",
			"./",
			remoteLocation,
			"--delete",
		}
	case Download:
		args = []string{
			"-ahr",
			"--perms",
			remoteLocation,
			".",
		}
	}

	cmd := exec.Command("rsync", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v failed to run rsync\n%s", err, string(out))
	}

	return nil
}
