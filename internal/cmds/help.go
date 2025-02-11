package cmds

import (
	"fmt"
	"os"
	"text/template"
	"time"
)

const help = `ritcs ...

version: v0.1.0

© mattglei.ch {{.Year}} [https://github.com/gleich/ritcs]

commands:
	setup     configure ritcs
`

type helpTemplate struct {
	Year int
}

func OutputHelp() error {
	tmpl, err := template.New("help").Parse(help)
	if err != nil {
		return fmt.Errorf("%v failed to parse template", err)
	}
	helpData := helpTemplate{Year: time.Now().Year()}
	err = tmpl.Execute(os.Stdout, helpData)
	if err != nil {
		return fmt.Errorf("%v failed to execute template", err)
	}
	return nil
}
