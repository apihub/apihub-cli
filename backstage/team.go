package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	httpErr "github.com/backstage/backstage/errors"
	"github.com/codegangsta/cli"
	. "github.com/mrvdot/golang-utils"
)

type Team struct {
	Name   string
	Alias  string
	client *Client
}

func (t *Team) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "team-create",
			Usage:       "team-create --name <name>",
			Description: "Creates a new team.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "Name of the team",
				},
			},
			Action: func(c *cli.Context) {
				defer RecoverStrategy("team-create")()
				team := &Team{
					Name:   c.String("name"),
					client: NewClient(&http.Client{}),
				}
				result := team.save()
				fmt.Println(result)
			},
		},
		{
			Name:        "team-info",
			Usage:       "team-info --name <name>",
			Description: "Retrieves team info.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "Name of the team",
				},
			},
			Action: func(c *cli.Context) {
				defer RecoverStrategy("team-info")()
				team := &Team{
					Name:   c.String("name"),
					client: NewClient(&http.Client{}),
				}
				result := team.info()
				fmt.Println(result)
			},
		},
		{
			Name:        "team-remove",
			Usage:       "team-remove --name <name>",
			Description: "Remove an existing team.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "Name of the team",
				},
			},
			Action: func(c *cli.Context) {
				context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
				if Confirm(context, "Are you sure you want to remove this team?") != true {
					fmt.Println(ErrCommandCancelled)
				} else {
					defer RecoverStrategy("team-remove")()
					team := &Team{
						Name:   c.String("name"),
						client: NewClient(&http.Client{}),
					}
					result := team.remove()
					fmt.Println(result)
				}
			},
		},
	}
}

func (t *Team) save() string {
	url, err := GetURL("/api/teams")
	if err != nil {
		return err.Error()
	}
	teamJson, err := json.Marshal(t)
	if err != nil {
		return err.Error()
	}
	b := bytes.NewBufferString(string(teamJson))
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err.Error()
	}

	response, err := t.client.Do(req)
	if err != nil {
		httpEr := err.(*httpErr.HTTPError)
		return httpEr.Message
	}

	if response.StatusCode == http.StatusCreated {
		return "Team created successfully."
	}
	return ErrBadRequest.Error()
}

func (t *Team) info() string {
	t.Alias = GenerateSlug(t.Name)
	url, err := GetURL("/api/teams" + "/" + t.Alias)
	if err != nil {
		return err.Error()
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err.Error()
	}

	response, err := t.client.Do(req)
	if err != nil {
		httpEr := err.(*httpErr.HTTPError)
		return httpEr.Message
	}

	if response.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		return string(body)
	}
	return ErrBadRequest.Error()
}

func (t *Team) remove() string {
	t.Alias = GenerateSlug(t.Name)
	url, err := GetURL("/api/teams" + "/" + t.Alias)
	if err != nil {
		return err.Error()
	}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err.Error()
	}

	response, err := t.client.Do(req)
	if err != nil {
		httpEr := err.(*httpErr.HTTPError)
		return httpEr.Message
	}

	if response.StatusCode == http.StatusOK {
		return "Team removed successfully."
	}
	return ErrBadRequest.Error()
}
