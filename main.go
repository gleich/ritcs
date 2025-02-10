package main

import (
	"time"

	"github.com/alecthomas/kong"
	"pkg.mattglei.ch/ritcs/internal/cmds"
	"pkg.mattglei.ch/timber"
)

func main() {
	timber.SetTimezone(time.Local)
	timber.SetTimeFormat("03:04:05")

	cli := struct {
		Setup struct{} `cmd:"" help:"configure ritcs with an interactive prompt"`
		Run   struct {
			Arguments []string `arg:"" help:"arguments to pass into the given command"`
			NoCopy    bool     `help:"do not copy the files from the output directory to the current directory"`
		} `cmd:"" help:"run a command on the RIT CS servers"`
	}{}

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
	case "run <arguments>":
		err := cmds.Run(cli.Run.Arguments)
		if err != nil {
			timber.Fatal(err, "failed to run \"run\" command")
		}
	}
}
