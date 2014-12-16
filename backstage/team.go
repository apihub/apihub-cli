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
			Description: "Create a team.",
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
			Description: "Return team info and lists of members and services.",
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
				table, err := team.info()
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
			Name:        "team-list",
			Usage:       "team-list",
			Description: "Return a list of all teams.",
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
			Description: "Delete a team.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "alias, a",
					Value: "",
					Usage: "Team alias",
				},
			},
			Action: func(c *cli.Context) {
				context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
				if Confirm(context, "Are you sure you want to delete this team? This action cannot be undone.") != true {
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
			Description: "Add a user to a team.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "team, t",
					Value: "",
					Usage: "Name of the team",
				},
				cli.StringFlag{
					Name:  "email, e",
					Value: "",
					Usage: "User's email",
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
			Description: "Remove a user from a team.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "team, t",
					Value: "",
					Usage: "Name of the team",
				},
				cli.StringFlag{
					Name:  "email, e",
					Value: "",
					Usage: "User's email",
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
		return "Your team has been created."
	}
	return ErrBadRequest.Error()
}

func (t *Team) info() (*Table, error) {
	path := "/api/teams/" + t.Alias
	var team Team
	response, err := t.client.MakeGet(path, &team)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusOK {
		fmt.Println("Name: " + team.Name)
		fmt.Println("Alias: " + team.Alias)
		fmt.Println("Owner: " + team.Owner + "\n")
		table := &Table{
			Content: [][]string{},
			Header:  []string{"Team Members"},
		}
		for _, member := range team.Users {
			line := []string{}
			line = append(line, member)
			table.Content = append(table.Content, line)
		}
		return table, nil
	}
	return nil, ErrBadRequest
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
		return "Your team has been deleted."
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
	return "Sorry, the user was not found."
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
	return "It's not allowed to remove the owner from its team."
}

func (t *Team) containsUserByEmail(email string) bool {
	for _, u := range t.Users {
		if u == email {
			return true
		}
	}
	return false
}
