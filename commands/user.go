package commands

import (
	"fmt"
	"os"

	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/backstage/backstage-cli/backstage"
	"github.com/codegangsta/cli"
)

type User struct {
	Service *backstage.UserService
}

func (cmd *User) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "user-create",
			Usage:       "user-create --name <name> --email <email> --username <username>",
			Description: "Create a user account.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "name, n", Value: "", Usage: "The user's real life name"},
				cli.StringFlag{Name: "email, e", Value: "", Usage: "User's email"},
				cli.StringFlag{Name: "username, u", Value: "", Usage: "Username is a unique variation on your name"},
			},
			Action: cmd.userCreate,
		},

		{
			Name:        "user-remove",
			Usage:       "user-remove",
			Description: "Delete a user account.",
			Action:      cmd.userDelete,
		},
	}
}

func (cmd *User) userCreate(c *cli.Context) {
	fmt.Println("Password (typing will be hidden):")
	password, _ := terminal.ReadPassword(int(os.Stdin.Fd()))

	_, err := cmd.Service.Create(c.String("name"), c.String("username"), c.String("email"), string(password))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Your account has been created.")
	}
}

func (cmd *User) userDelete(c *cli.Context) {
	defer RecoverStrategy("user-remove")()

	context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
	if Confirm(context, "Are you sure you want to delete your account? If deleted, you can't restore it.") != true {
		fmt.Println(backstage.ErrCommandCancelled.Error())
		return
	}

	err := cmd.Service.Delete()

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Your account has been deleted.")
	}
}
