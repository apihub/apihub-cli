package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate
	app := cli.NewApp()
	app.Name = "backstage"
	app.Usage = "An open source solution for publishing APIs."
	app.Version = "0.0.1"
	app.HideHelp = true

	m := NewManager(app)
	m.Register(&Auth{})
	m.Register(&Target{})

	app.Run(os.Args)
}
