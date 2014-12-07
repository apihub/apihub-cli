package main

import (
	"errors"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	cli.AppHelpTemplate = AppHelpTemplate
	cmd := cli.NewApp()
	cmd.Name = "backstage"
	cmd.Usage = "An open source solution for publishing APIs."
	cmd.Version = "0.0.1"
	cmd.HideHelp = true
	cmd.Commands = []cli.Command{
		{
			Name:        "login",
			Usage:       "login <email>",
			Description: "Sign in with your Backstage credentials to continue.",
			Before: func(c *cli.Context) error {
				context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
				if Confirm(context, "Are you sure you want to login?") != true {
					return errors.New("")
				}
				return nil
			},
			Action: func(c *cli.Context) {
				args := c.Args().Tail()
				println("login: ", args[0])
			},
		},
	}
	cmd.Run(os.Args)
}
