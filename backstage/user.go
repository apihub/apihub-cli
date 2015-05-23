package main

import (
	"fmt"
	"net/http"
	"os"

	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/codegangsta/cli"
)

type User struct {
	Name                 string      `json:"name,omitempty"`
	Email                string      `json:"email,omitempty"`
	Username             string      `json:"username,omitempty"`
	Password             string      `json:"password,omitempty"`
	NewPassword          string      `json:"new_password,omitempty"`
	ConfirmationPassword string      `json:"confirmation_password,omitempty"`
	client               *HTTPClient `json:"-"`
}

func (u *User) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "user-create",
			Usage:       "user-create --name <name> --email <email> --username <username>",
			Description: "Create a user account.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "The user's real life name",
				},
				cli.StringFlag{
					Name:  "email, e",
					Value: "",
					Usage: "User's email",
				},
				cli.StringFlag{
					Name:  "username, u",
					Value: "",
					Usage: "Username is a unique variation on your name",
				},
			},
			Action: func(c *cli.Context) {
				fmt.Println("Password (typing will be hidden):")
				password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					fmt.Println(err.Error())
				}
				user := &User{
					Name:     c.String("name"),
					Email:    c.String("email"),
					Username: c.String("username"),
					Password: string(password),
					client:   NewHTTPClient(&http.Client{}),
				}
				result := user.save()
				fmt.Println(result)
			},
		},
		{
			Name:        "user-remove",
			Usage:       "user-remove",
			Description: "Delete a user account.",
			Action: func(c *cli.Context) {
				context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
				if Confirm(context, "Are you sure you want to delete your account? If deleted, you can't restore it.") != true {
					fmt.Println(ErrCommandCancelled.Error())
					return
				}

				defer RecoverStrategy("user-remove")()
				user := &User{
					client: NewHTTPClient(&http.Client{}),
				}
				result := user.remove()
				fmt.Println(result)
			},
		},
	}
}

func (u *User) save() string {
	path := "/api/users"
	user := &User{}
	response, err := u.client.MakePost(path, u, user)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusCreated {
		return "Your account has been created."
	}
	return ErrBadRequest.Error()
}

func (u *User) remove() string {
	path := "/api/users"
	user := User{}
	response, err := u.client.MakeDelete(path, nil, &user)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusOK {
		DeleteToken()
		return "Your account has been deleted."
	}
	return ErrBadRequest.Error()
}
