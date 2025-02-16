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
	timber.Timezone(time.Local)
	timber.TimeFormat("03:04:05")

	if len(os.Args) < 2 {
		cmds.OutputHelp()
		fmt.Println()
		timber.FatalMsg("please provide command to run")
	}

	var (
		skipDownloadFlag bool
		silentFlag       bool
	)
	args := os.Args[1:]
	for len(args) > 0 && (args[0] == "--skip-download" || args[0] == "--silent") {
		switch args[0] {
		case "--skip-download":
			skipDownloadFlag = true
		case "--silent":
			silentFlag = true
		}
		args = args[1:]
	}

	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		cmds.OutputHelp()
		return
	}

	if len(args) > 0 {
		switch strings.ToLower(strings.TrimSpace(args[0])) {
		case "setup":
			cmds.Setup()
			return
		case "uninstall":
			cmds.Uninstall()
			return
		}
	}

	if err := conf.Load(); err != nil {
		timber.Fatal(err, "failed to load configuration file")
	}
	if skipDownloadFlag {
		conf.Config.SkipDownload = true
	}
	if silentFlag {
		conf.Config.Silent = true
	}

	cmds.Run(args)
}
