package commands

import (
	"fmt"
	"os"

	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/apihub/apihub-cli/maestro"
	"github.com/codegangsta/cli"
)

type Auth struct {
	Service *apihub.AuthService
}

func (cmd *Auth) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "login",
			Usage:       "login <email>",
			Description: "Login in with your ApiHub credentials.",
			Action:      cmd.login,
		},

		{
			Name:        "logout",
			Usage:       "logout",
			Description: "Clear local credentials.",
			Action:      cmd.logout,
		},

		{
			Name:        "change-password",
			Usage:       "change-password <email>",
			Description: "Change the password of the user provided.",
			Action:      cmd.changePassword,
		},
	}
}

func (cmd *Auth) login(c *cli.Context) {
	email := c.Args().First()
	fmt.Println("Password (typing will be hidden):")
	password, _ := terminal.ReadPassword(int(os.Stdin.Fd()))

	_, err := cmd.Service.Login(email, string(password))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Authentication successful.")
	}
}

func (cmd *Auth) logout(c *cli.Context) {
	err := cmd.Service.Logout()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("You have successfully logged out.")
	}
}

func (cmd *Auth) changePassword(c *cli.Context) {
	email := c.Args().First()
	fmt.Println("Current password (typing will be hidden):")
	password, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println("New password (typing will be hidden):")
	newPassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println("Confirme new password (typing will be hidden):")
	confirmationPassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))

	err := cmd.Service.ChangePassword(email, string(password), string(newPassword), string(confirmationPassword))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("You password has been changed.")
	}
}
