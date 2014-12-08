package main

import (
	"fmt"
	"os"

	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/codegangsta/cli"
)

type Auth struct{}

func (a *Auth) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "login",
			Usage:       "login <email>",
			Description: "Sign in with your Backstage credentials to continue.",
			Action: func(c *cli.Context) {
				password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					fmt.Print(err)
					return
				}
				fmt.Printf("password %+v\n", string(password))
			},
		},
	}
}
