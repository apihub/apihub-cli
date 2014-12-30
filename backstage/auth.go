package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"syscall"

	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/codegangsta/cli"
)

var TokenFileName  = joinHomePath(".backstage_token")

type Auth struct {
	client *HTTPClient
}

func (a *Auth) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "login",
			Usage:       "login <email>\n   Password (typing will be hidden):\n   Authentication successful.",
			Description: "Login in with your Backstage credentials.",
			Action: func(c *cli.Context) {
				email := c.Args().First()
				fmt.Println("Password (typing will be hidden):")
				password, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
				auth := &Auth{client: NewHTTPClient(&http.Client{})}
				result := auth.Login(email, string(password))
				fmt.Println(result)
			},
		},
		{
			Name:        "logout",
			Usage:       "logout\n   You have successfully logged out.",
			Description: "Clear local credentials.",
			Action: func(c *cli.Context) {
				auth := &Auth{}
				result := auth.Logout()
				fmt.Println(result)
			},
		},
	}
}

func (a *Auth) Login(email, password string) string {
	path := "/api/login"
	user := &User{
		Email:    email,
		Password: password,
	}
	token := map[string]interface{}{}
	_, err := a.client.MakePost(path, user, &token)
	if err != nil {
		return err.Error()
	}

	if err := writeToken(token["token_type"].(string) + " " + token["access_token"].(string)); err != nil {
		return err.Error()
	}
	return "Authentication successful."
}

func (a *Auth) Logout() string {
	filesystem().Remove(joinHomePath(".backstage_token"))
	return "You have successfully logged out."
}

func writeToken(token string) error {
	tokenFile, err := filesystem().OpenFile(TokenFileName, syscall.O_RDWR|syscall.O_CREAT|syscall.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer tokenFile.Close()
	tokenFile.WriteString(token)
	return nil
}

func ReadToken() (string, error) {
	tokenFile, err := filesystem().OpenFile(TokenFileName, syscall.O_RDWR, 0600)
	if err != nil {
		return "", ErrLoginRequired
	}
	defer tokenFile.Close()
	data, _ := ioutil.ReadAll(tokenFile)
	return string(data), nil
}

func DeleteToken() error {
	err := filesystem().Remove(TokenFileName)
	return err
}
