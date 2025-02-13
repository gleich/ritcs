package cmds

import (
	"fmt"
	"time"
)

func OutputHelp() {
	fmt.Printf(`ritcs ...

version: v1.0.0

© mattglei.ch %d [https://github.com/gleich/ritcs]

commands:
	setup      configure ritcs
	uninstall  remove all ritcs files locally and from remote machine
`, time.Now().Year())
}
