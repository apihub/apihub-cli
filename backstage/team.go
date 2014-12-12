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
	"gopkg.in/mgo.v2/bson"
)

type Team struct {
	Id     bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty""`
	Name   string        `json:"name"`
	Alias  string        `json:"alias"`
	Users  []string      `json:"users"`
	Owner  string        `json:"owner"`
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
			Usage:       "team-info --alias <alias>",
			Description: "Retrieves team info.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "alias, a",
					Value: "",
					Usage: "Team alias",
				},
			},
			Action: func(c *cli.Context) {
				defer RecoverStrategy("team-info")()
				team := &Team{
					Alias:  c.String("alias"),
					client: NewClient(&http.Client{}),
				}
				result := team.info()
				fmt.Println(result)
			},
		},
		{
			Name:        "team-remove",
			Usage:       "team-remove --alias <alias>",
			Description: "Remove an existing team.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "alias, a",
					Value: "",
					Usage: "Team alias",
				},
			},
			Action: func(c *cli.Context) {
				context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
				if Confirm(context, "Are you sure you want to remove this team?") != true {
					fmt.Println(ErrCommandCancelled)
				} else {
					defer RecoverStrategy("team-remove")()
					team := &Team{
						Alias:  c.String("alias"),
						client: NewClient(&http.Client{}),
					}
					result := team.remove()
					fmt.Println(result)
				}
			},
		},
		{
			Name:        "team-user-add",
			Usage:       "team-user-add --team <team-alias> --email <email>",
			Description: "Adds a user in a team.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "team, t",
					Value: "",
					Usage: "Name of the team",
				},
				cli.StringFlag{
					Name:  "email, e",
					Value: "",
					Usage: "User email",
				},
			},
			Action: func(c *cli.Context) {
				defer RecoverStrategy("team-user-add")()
				team := &Team{
					Alias:  c.String("team"),
					client: NewClient(&http.Client{}),
				}
				result := team.addUser(c.String("email"))
				fmt.Println(result)
			},
		},
		{
			Name:        "team-user-remove",
			Usage:       "team-user-remove --team <team-alias> --email <email>",
			Description: "Removes a user from a team.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "team, t",
					Value: "",
					Usage: "Name of the team",
				},
				cli.StringFlag{
					Name:  "email, e",
					Value: "",
					Usage: "User email",
				},
			},
			Action: func(c *cli.Context) {
				defer RecoverStrategy("team-user-remove")()
				team := &Team{
					Alias:  c.String("team"),
					client: NewClient(&http.Client{}),
				}
				result := team.removeUser(c.String("email"))
				fmt.Println(result)
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

func (t *Team) addUser(user string) string {
	url, err := GetURL("/api/teams/" + t.Alias + "/users")
	if err != nil {
		return err.Error()
	}
	t.Users = append(t.Users, user)
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

	var team = &Team{}
	parseBody(response.Body, &team)
	if response.StatusCode == http.StatusCreated && team.containsEmail(t.Users[0]) {
		return "User `" + t.Users[0] + "` added successfully to team `" + t.Alias + "`."
	}
	return "User not found! Please check if the email provided is a valid user in the server."
}

func (t *Team) removeUser(user string) string {
	url, err := GetURL("/api/teams/" + t.Alias + "/users")
	if err != nil {
		return err.Error()
	}
	t.Users = append(t.Users, user)
	teamJson, err := json.Marshal(t)
	if err != nil {
		return err.Error()
	}
	b := bytes.NewBufferString(string(teamJson))
	req, err := http.NewRequest("DELETE", url, b)
	if err != nil {
		return err.Error()
	}

	response, err := t.client.Do(req)
	if err != nil {
		httpEr := err.(*httpErr.HTTPError)
		return httpEr.Message
	}

	var team = &Team{}
	parseBody(response.Body, &team)
	if response.StatusCode == http.StatusOK && !team.containsEmail(t.Users[0]) {
		return "User `" + t.Users[0] + "` removed successfully to team `" + t.Alias + "`."
	}
	if team.Owner == user {
		return "You cannot remove the owner."
	}
	return "User not found! Please check if the email provided is a valid user in the server."
}

func (t *Team) containsEmail(email string) bool {
	for _, u := range t.Users {
		if u == email {
			return true
		}
	}
	return false
}
