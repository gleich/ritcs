package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"pkg.mattglei.ch/ritcs/internal/cmds"
	"pkg.mattglei.ch/ritcs/internal/conf"
	"pkg.mattglei.ch/timber"
)

func main() {
	timber.SetTimezone(time.Local)
	timber.SetTimeFormat("03:04:05")

	if len(os.Args) < 2 {
		err := cmds.OutputHelp()
		if err != nil {
			timber.Fatal(err, "failed to output help")
		}
		fmt.Println()
		timber.FatalMsg("please provide command to run")
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		err := cmds.OutputHelp()
		if err != nil {
			timber.Fatal(err, "failed to output help")
		}
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

	err = cmds.Run(args)
	if err != nil {
		timber.Fatal(err, "failed to run command")
	}
}
