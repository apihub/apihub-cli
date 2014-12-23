package main

import (
	"os"

	"github.com/codegangsta/cli"
)

var BackstageClientVersion = "0.0.6"

func main() {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate
	app := cli.NewApp()
	app.Name = "backstage"
	app.Usage = "An open source solution for publishing APIs."
	app.Version = BackstageClientVersion
	app.HideHelp = true

	m := NewManager(app)
	m.Register(&Auth{})
	m.Register(&Service{})
	m.Register(&Target{})
	m.Register(&Team{})
	m.Register(&User{})

	app.Run(os.Args)
}
