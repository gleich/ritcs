package cmds

import (
	"fmt"
	"time"
)

func OutputHelp() {
	fmt.Printf(`ritcs [flags] <command> [arguments...]

version: v1.0.1

© mattglei.ch %d [https://github.com/gleich/ritcs]

flags:
  --skip-download  skip downloading remote changes.
  --silent         run in silent mode.

commands:
  setup      configure ritcs.
  uninstall  remove all ritcs files locally and from remote machine.
`, time.Now().Year())
}
