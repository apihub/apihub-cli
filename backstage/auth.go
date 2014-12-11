package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"syscall"

	"code.google.com/p/go.crypto/ssh/terminal"
	httpErr "github.com/backstage/backstage/errors"
	"github.com/codegangsta/cli"
)

var (
	TokenFileName              = joinHomePath(".backstage_token")
	ErrFailedWrittingTokenFile = errors.New("Failed trying to write the token file.")
	ErrLoginRequired           = errors.New("You must log in.")
	ErrBadRequest              = errors.New("It was not possible to handle your request at this moment. Please ty again.")
)

type Auth struct {
	client *Client
}

func (a *Auth) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "login",
			Usage:       "login <email>",
			Description: "Sign in with your Backstage credentials to continue.",
			Action: func(c *cli.Context) {
				email := c.Args().First()
				fmt.Println("Please type your password:")
				password, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
				auth := &Auth{client: NewClient(&http.Client{})}
				result := auth.Login(email, string(password))
				fmt.Println(result)
			},
		},
		{
			Name:        "logout",
			Description: "Sign out from Backstage.",
			Action: func(c *cli.Context) {
				auth := &Auth{}
				resp := auth.Logout()
				fmt.Println(resp)
			},
		},
	}
}

func (a *Auth) Login(email, password string) string {
	url, err := GetURL("/api/login")
	if err != nil {
		return err.Error()
	}
	b := bytes.NewBufferString(`{"email":"` + email + `", "password":"` + password + `"}`)
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err.Error()
	}

	response, err := a.client.Do(req)
	if err != nil {
		httpEr := err.(*httpErr.HTTPError)
		return httpEr.Message
	}

	var token = map[string]interface{}{}
	parseBody(response.Body, &token)
	writeToken(token["token_type"].(string) + " " + token["token"].(string))
	return "Welcome! You've signed in successfully."
}

func (a *Auth) Logout() string {
	err := filesystem().Remove(joinHomePath(".backstage_token"))
	if err != nil {
		return "You are not signed in."
	}
	return "You have signed out successfully."
}

func writeToken(token string) error {
	tokenFile, err := filesystem().OpenFile(TokenFileName, syscall.O_RDWR|syscall.O_CREAT|syscall.O_TRUNC, 0600)
	defer tokenFile.Close()
	if err != nil {
		return err
	}
	tokenFile.WriteString(token)
	return nil
}

func ReadToken() (string, error) {
	tokenFile, err := filesystem().OpenFile(TokenFileName, syscall.O_RDWR, 0600)
	defer tokenFile.Close()
	if err != nil {
		return "", ErrLoginRequired
	}
	data, _ := ioutil.ReadAll(tokenFile)
	return string(data), nil
}

func DeleteToken() error {
	err := filesystem().Remove(TokenFileName)
	if err != nil {
		return err
	}
	return nil
}
