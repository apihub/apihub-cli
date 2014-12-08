package main

import (
	"fmt"
	"os"

	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/codegangsta/cli"
)

func main() {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate
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
			Action: func(c *cli.Context) {
				password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					fmt.Print(err)
					return
				}
				fmt.Printf("password %+v\n", string(password))
			},
		},
		{
			Name:        "target-add",
			Usage:       "target-add <label> <endpoint>",
			Description: "Adds a new target in the list of targets.",
			Action: func(c *cli.Context) {
				defer RecoverStrategy("target-add")()
				targets, err := LoadTargets()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				args := c.Args()
				label, endpoint := args[0], args[1]
				err = targets.add(label, endpoint)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println("Target added successfully!")
			},
		},
		{
			Name:        "target-list",
			Usage:       "",
			Description: "Adds a new target in the list of targets.",
			Action: func(c *cli.Context) {
				targets, err := LoadTargets()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(targets.list())
			},
		},
		{
			Name:        "target-remove",
			Usage:       "target-remove <label>",
			Description: "Remove a target from the list of targets.",
			Before: func(c *cli.Context) error {
				if c.Args().First() == "" {
					return ErrCommandCancelled
				}
				context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
				if Confirm(context, "Are you sure you want to remove this target?") != true {
					return ErrCommandCancelled
				}
				return nil
			},
			Action: func(c *cli.Context) {
				defer RecoverStrategy("target-remove")()
				targets, err := LoadTargets()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				label := c.Args()[1]
				err = targets.remove(label)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println("Target removed successfully!")
			},
		},
		{
			Name:        "target-set",
			Usage:       "target-set <label>",
			Description: "Set a target as default to be used.",
			Action: func(c *cli.Context) {
				defer RecoverStrategy("target-set")()
				targets, err := LoadTargets()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				label := c.Args().First()
				err = targets.setDefault(label)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println("You have a new target as default!")
			},
		},
	}
	cmd.Run(os.Args)
}
