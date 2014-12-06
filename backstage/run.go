package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "backstage"
	app.Usage = "An open source solution for publishing APIs."
	app.Version = "0.0.1"
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Name:  "login",
			Usage: "login <email>",
			Action: func(c *cli.Context) {
				println("login: ", c.Args().First())
			},
		},
	}
	app.Run(os.Args)
}
