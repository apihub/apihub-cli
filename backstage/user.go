package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/codegangsta/cli"
)

type User struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	client   *Client
}

func (u *User) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "user-create",
			Usage:       "user-create --name <name> --email <email> --username <username>",
			Description: "Creates a new account.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "Name of the user",
				},
				cli.StringFlag{
					Name:  "email, e",
					Value: "",
					Usage: "Email",
				},
				cli.StringFlag{
					Name:  "username, u",
					Value: "",
					Usage: "Username",
				},
			},
			Action: func(c *cli.Context) {
				defer RecoverStrategy("user-create")()
				fmt.Println("Please type your password:")
				password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					fmt.Println(err.Error())
				}
				user := &User{
					Name:     c.String("name"),
					Email:    c.String("email"),
					Username: c.String("username"),
					Password: string(password),
					client:   NewClient(&http.Client{}),
				}
				result := user.save()
				fmt.Println(result)
			},
		},
		{
			Name:        "user-remove",
			Usage:       "user-remove",
			Description: "Removes an account.",
			Before: func(c *cli.Context) error {
				context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
				if Confirm(context, "Are you sure you want to remove your user?(Everything will be lost!)") != true {
					return ErrCommandCancelled
				}
				return LoginRequired()
			},
			Action: func(c *cli.Context) {
				defer RecoverStrategy("user-remove")()
				user := &User{
					client: NewClient(&http.Client{}),
				}
				result := user.remove()
				fmt.Println(result)
			},
		},
	}
}

func (u *User) save() string {
	path := "/api/users"
	payload, err := json.Marshal(u)
	if err != nil {
		return err.Error()
	}
	user := User{}
	response, err := u.client.MakePost(path, string(payload), &user)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusCreated {
		return "User created successfully."
	}
	return ErrBadRequest.Error()
}

func (u *User) remove() string {
	path := "/api/users"
	user := User{}
	response, err := u.client.MakeDelete(path, "", &user)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusOK {
		DeleteToken()
		return "User removed successfully."
	}
	return ErrBadRequest.Error()
}
