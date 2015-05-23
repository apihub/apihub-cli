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

var TokenFileName = joinHomePath(".backstage_token")

type Auth struct {
	client *HTTPClient
}

func (a *Auth) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "login",
			Usage:       "login <email>",
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
			Usage:       "logout",
			Description: "Clear local credentials.",
			Action: func(c *cli.Context) {
				auth := &Auth{client: NewHTTPClient(&http.Client{})}
				result := auth.Logout()
				fmt.Println(result)
			},
		},
		{
			Name:        "change-password",
			Usage:       "change-password <email>",
			Description: "Change the password of the user provided.",
			Action: func(c *cli.Context) {
				email := c.Args().First()
				fmt.Println("Current password (typing will be hidden):")
				password, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
				fmt.Println("New password (typing will be hidden):")
				newPassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
				fmt.Println("Confirme new password (typing will be hidden):")
				confirmationPassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
				auth := &Auth{client: NewHTTPClient(&http.Client{})}
				result := auth.ChangePassword(email, string(password), string(newPassword), string(confirmationPassword))
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
	path := "/api/logout"
	response, err := a.client.MakeDelete(path, nil, nil)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusNoContent {
		filesystem().Remove(joinHomePath(".backstage_token"))
		return "You have successfully logged out."
	}
	return ErrBadRequest.Error()
}

func (a *Auth) ChangePassword(email, password, newPassword, confirmationPassword string) string {
	path := "/api/password"
	user := &User{
		Email:                email,
		Password:             password,
		NewPassword:          newPassword,
		ConfirmationPassword: confirmationPassword,
	}
	response, err := a.client.MakePut(path, user, nil)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusNoContent {
		return "You password has been changed."
	}
	return "It was not possible to change your password. Please try again."
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
