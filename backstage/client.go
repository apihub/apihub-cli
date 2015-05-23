package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
)

type Client struct {
	Id          string `json:"id,omitempty"`
	Secret      string `json:"secret,omitempty"`
	Name        string `json:"name,omitempty"`
	RedirectUri string `json:"redirect_uri,omitempty"`
	Team        string `json:"team,omitempty"`
	client      *HTTPClient
}

func (c *Client) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "client-create",
			Usage:       "client-create --team <team> --client_id <client_id> --name <name> --redirect_uri <redirect_uri>",
			Description: "Create a new client.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "client_id, i",
					Value: "",
					Usage: "Client id (used by oAuth 2.0)",
				},
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "Client name",
				},
				cli.StringFlag{
					Name:  "redirect_uri, r",
					Value: "",
					Usage: "Client Redirect Uri (used by oAuth 2.0)",
				},
				cli.StringFlag{
					Name:  "team, t",
					Value: "",
					Usage: "Team",
				},
			},
			Action: func(c *cli.Context) {
				defer RecoverStrategy("client-create")()
				client := &Client{
					Id:          c.String("client_id"),
					Name:        c.String("name"),
					RedirectUri: c.String("redirect_uri"),
					Team:        c.String("team"),
					client:      NewHTTPClient(&http.Client{}),
				}
				result := client.create()
				fmt.Println(result)
			},
		},
		{
			Name:        "client-update",
			Usage:       "client-update --team <team> --client_id <client_id> --name <name> --redirect_uri <redirect_uri>",
			Description: "Update an existing client.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "client_id, i",
					Value: "",
					Usage: "Client id (used by oAuth 2.0)",
				},
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "Client name",
				},
				cli.StringFlag{
					Name:  "redirect_uri, r",
					Value: "",
					Usage: "Client Redirect Uri (used by oAuth 2.0)",
				},
				cli.StringFlag{
					Name:  "team, t",
					Value: "",
					Usage: "Team",
				},
			},
			Action: func(c *cli.Context) {
				defer RecoverStrategy("client-update")()
				client := &Client{
					Id:          c.String("client_id"),
					Name:        c.String("name"),
					RedirectUri: c.String("redirect_uri"),
					Team:        c.String("team"),
					client:      NewHTTPClient(&http.Client{}),
				}
				result := client.update()
				fmt.Println(result)
			},
		},
		{
			Name:        "client-info",
			Usage:       "client-info --client_id <client_id>",
			Description: "Retrieve client information.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "client_id, i",
					Value: "",
					Usage: "Client Id",
				},
			},
			Action: func(c *cli.Context) {
				defer RecoverStrategy("client-info")()
				client := &Client{
					Id:     c.String("client_id"),
					client: NewHTTPClient(&http.Client{}),
				}
				table, err := client.info()
				if err != "" {
					fmt.Println(err)
				}

				if table != nil {
					context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
					table.Render(context)
					fmt.Println("\n")
				}
			},
		},
		{
			Name:        "client-remove",
			Usage:       "client-remove --client_id <client_id> --team <team>\n   The client `<name>` has been deleted.",
			Description: "Remove an existing client.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "client_id, i",
					Value: "",
					Usage: "Client Id",
				},
				cli.StringFlag{
					Name:  "team, t",
					Value: "",
					Usage: "Team",
				},
			},
			Action: func(c *cli.Context) {
				context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
				if Confirm(context, "Are you sure you want to delete this client? This action cannot be undone.") != true {
					fmt.Println(ErrCommandCancelled)
				} else {
					defer RecoverStrategy("client-remove")()
					client := &Client{
						Id:     c.String("client_id"),
						Team:   c.String("team"),
						client: NewHTTPClient(&http.Client{}),
					}
					result := client.remove()
					fmt.Println(result)
				}
			},
		},
	}
}

func (c *Client) create() string {
	path := fmt.Sprintf("/api/teams/%s/clients", c.Team)
	client := &Client{}
	response, err := c.client.MakePost(path, c, client)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusCreated {
		return "Your new client has been created."
	}
	panic("The client was not found for the team provided.")
}

func (c *Client) update() string {
	path := fmt.Sprintf("/api/teams/%s/clients/%s", c.Team, c.Id)
	client := &Client{}
	response, err := c.client.MakePut(path, c, client)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusOK {
		return "Your client has been updated."
	}
	panic("The client was not found for the team provided.")
}

func (c *Client) info() (*Table, string) {
	path := fmt.Sprintf("/api/teams/clients/%s", c.Id)
	var client Client
	response, err := c.client.MakeGet(path, &client)
	if err != nil {
		return nil, err.Error()
	}

	if response.StatusCode == http.StatusOK {
		fmt.Println("Team Name: " + client.Team)
		fmt.Println("")
		clientTable := &Table{
			Title:   "Clients Details:",
			Content: [][]string{},
			Header:  []string{"Name", "Redirect Uri", "Id", "Secret"},
		}
		line := []string{}
		line = append(line, client.Name)
		line = append(line, client.RedirectUri)
		line = append(line, client.Id)
		line = append(line, client.Secret)
		clientTable.Content = append(clientTable.Content, line)
		return clientTable, ""
	}
	panic("The client was not found for the team provided.")
}

func (c *Client) remove() string {
	path := fmt.Sprintf("/api/teams/%s/clients/%s", c.Team, c.Id)
	client := &Client{}
	response, err := c.client.MakeDelete(path, nil, client)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusOK {
		return "The client `" + c.Id + "` has been deleted."
	}
	panic("The client was not found for the team provided.")
}
