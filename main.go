package main

import (
	"os"
	"time"

	"github.com/alecthomas/kong"
	"pkg.mattglei.ch/ritcs/internal/cmds"
	"pkg.mattglei.ch/ritcs/internal/conf"
	"pkg.mattglei.ch/ritcs/internal/remote"
	"pkg.mattglei.ch/timber"
)

var CLI struct {
	Setup struct{} `cmd:"" help:"configure ritcs with an interactive prompt"`
}

func main() {
	timber.SetTimezone(time.Local)
	timber.SetTimeFormat("03:04:05")

	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "setup":
		err := cmds.Setup()
		if err != nil {
			timber.Fatal(err, "failed to setup user")
		}
		return
	default:
	}

	if len(os.Args) <= 1 {
		timber.FatalMsg("please provide arguments to get command")
	}

	conf, err := conf.Load()
	if err != nil {
		timber.Fatal(err, "failed to load configuration file")
	}

	sshClient, sftpClient, err := remote.EstablishConnection(conf)
	if err != nil {
		timber.Fatal(err, "failed to establish connection")
	}
	defer sftpClient.Close()
	defer sshClient.Close()

	tempDir, err := remote.CreateTempDir(conf, sftpClient)
	if err != nil {
		timber.Fatal(err, "failed to create temporary directory on server")
	}

	err = remote.RunGet(sshClient, tempDir)
	if err != nil {
		timber.Fatal(err, "failed to run get command")
	}

	err = remote.CopyFilesFromDir(sftpClient, tempDir)
	if err != nil {
		timber.Fatal(err, "failed to copy over files")
	}

	err = remote.RemoveTempDir(sftpClient, tempDir)
	if err != nil {
		timber.Fatal(err, "failed to remove temporary directory")
	}
}
