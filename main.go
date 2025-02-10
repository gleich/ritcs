package main

import (
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
		timber.FatalMsg("please provide command to run")
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
