package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"pkg.mattglei.ch/ritcs/internal/cmds"
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

	err := cmds.Run(os.Args[1:])
	if err != nil {
		timber.Fatal(err, "failed to run command")
	}
}
