package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

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
			Name:        "team-list",
			Usage:       "team-list",
			Description: "Retrieves all your teams.",
			Action: func(c *cli.Context) {
				defer RecoverStrategy("team-list")()
				team := &Team{
					client: NewClient(&http.Client{}),
				}
				table, err := team.list()
				if table != nil {
					context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
					table.Render(context)
					return
				}
				if err != nil {
					fmt.Println(err.Error())
					return
				}
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
	path := "/api/teams"
	team := &Team{}
	response, err := t.client.MakePost(path, t, team)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusCreated {
		return "Team created successfully."
	}
	return ErrBadRequest.Error()
}

func (t *Team) info() string {
	path := "/api/teams/" + t.Alias
	team := &Team{}
	response, err := t.client.MakeGet(path, team)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusOK {
		return team.Name
	}
	return ErrBadRequest.Error()
}

func (t *Team) list() (*Table, error) {
	path := "/api/teams"
	var teams []map[string]string
	_, err := t.client.MakeGet(path, &teams)
	if err != nil {
		return nil, err
	}

	if len(teams) > 0 {
		table := &Table{
			Content: [][]string{},
			Header:  []string{"Team Name", "Alias", "Owner"},
		}
		for _, team := range teams {
			line := []string{}
			line = append(line, team["name"], team["alias"], team["owner"])
			table.Content = append(table.Content, line)
		}
		return table, nil
	}
	return nil, errors.New("You have no teams.")
}

func (t *Team) remove() string {
	path := "/api/teams/" + t.Alias
	team := &Team{}
	response, err := t.client.MakeDelete(path, nil, team)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusOK {
		return "Team removed successfully."
	}
	return ErrBadRequest.Error()
}

func (t *Team) addUser(email string) string {
	path := "/api/teams/" + t.Alias + "/users"
	var team = &Team{}
	response, err := t.client.MakePost(path, t, team)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusCreated && team.containsUserByEmail(email) {
		return "User `" + email + "` added successfully to team `" + t.Alias + "`."
	}
	return "User not found! Please check if the email provided is a valid user in the server."
}

func (t *Team) removeUser(email string) string {
	path := "/api/teams/" + t.Alias + "/users"
	t.Users = append(t.Users, email)
	team := &Team{}
	response, err := t.client.MakeDelete(path, t, team)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusOK && !team.containsUserByEmail(email) {
		return "User `" + email + "` removed successfully to team `" + t.Alias + "`."
	}
	return "You cannot remove the owner."
}

func (t *Team) containsUserByEmail(email string) bool {
	for _, u := range t.Users {
		if u == email {
			return true
		}
	}
	return false
}
