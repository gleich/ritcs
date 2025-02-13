package cmds

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"go.mattglei.ch/ritcs/internal/conf"
	"go.mattglei.ch/ritcs/internal/remote"
	"go.mattglei.ch/timber"
)

func Uninstall() {
	sshClient, err := remote.EstablishConnection()
	if err != nil {
		timber.Fatal(err, "failed to establish connection to host machine")
	}
	defer sshClient.Close()

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		timber.Fatal(err, "failed to create sftp client")
	}
	defer sftpClient.Close()

	timber.Info("removing .ritcs directory from", conf.Config.Host)
	err = sftpClient.RemoveAll(remote.RemoteRITCSDirectory())
	if err != nil {
		timber.Fatal(err, "failed to remove remote .ritcs directory")
	}
	timber.Done("removed .ritcs directory from", conf.Config.Host)

	fmt.Println()
	timber.Info("removing local ritcs config")
	confPath, err := conf.Path()
	if err != nil {
		timber.Fatal(err, "failed to get configuration path")
	}
	confDir := filepath.Dir(confPath)
	err = os.RemoveAll(confDir)
	if err != nil {
		timber.Fatal(err, "failed to remove", confDir)
	}
	timber.Done("removed local ritcs config")
}
