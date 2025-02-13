package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.mattglei.ch/ritcs/internal/cmds"
	"go.mattglei.ch/ritcs/internal/conf"
	"go.mattglei.ch/timber"
)

func main() {
	timber.SetTimezone(time.Local)
	timber.SetTimeFormat("03:04:05")

	if len(os.Args) < 2 {
		cmds.OutputHelp()
		fmt.Println()
		timber.FatalMsg("please provide command to run")
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		cmds.OutputHelp()
		return
	}

	if strings.ToLower(strings.Trim(os.Args[1], " ")) == "setup" {
		err := cmds.Setup()
		if err != nil {
			timber.Fatal(err, "failed to setup user")
		}
		return
	}

	err := conf.Load()
	if err != nil {
		timber.Fatal(err, "failed to load configuration file")
	}

	args := os.Args[1:]
	if os.Args[1] == "--skip-download" {
		conf.Config.SkipDownload = true
		args = os.Args[2:]
	} else if os.Args[1] == "--silent" {
		conf.Config.Silent = true
		args = os.Args[2:]
	}

	if args[0] == "uninstall" {
		cmds.Uninstall()
	} else {
		cmds.Run(args)
	}
}
