package main

import (
	"time"

	"github.com/alecthomas/kong"
	"pkg.mattglei.ch/ritcs/internal/cmds"
	"pkg.mattglei.ch/timber"
)

var cli struct {
	Setup struct{} `cmd:"" help:"configure ritcs with an interactive prompt"`
	Get   struct {
		Arguments []string `arg:"" help:"arguments to pass to the get command"`
	} `cmd:"" help:"get files using the \"get\" command"`
}

func main() {
	timber.SetTimezone(time.Local)
	timber.SetTimeFormat("03:04:05")

	ctx := kong.Parse(
		&cli,
		kong.Name("ritcs"),
		kong.Description("ritcs is a tool to interact with the RIT CS servers"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
	)

	switch ctx.Command() {
	case "setup":
		err := cmds.Setup()
		if err != nil {
			timber.Fatal(err, "failed to setup user")
		}
		return
	case "get <arguments>":
		err := cmds.Get(cli.Get.Arguments)
		if err != nil {
			timber.Fatal(err, "failed to run get command")
		}
	default:
	}
}
