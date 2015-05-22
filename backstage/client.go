package main

import (
  "fmt"
  "net/http"
  "os"

  "github.com/codegangsta/cli"
)

type Client struct {
  Id          string `json:"id"`
  Name        string `json:"name"`
  RedirectUri string `json:"redirect_uri"`
  Team        string `json:"team"`
  client          *HTTPClient
}

func (c *Client) GetCommands() []cli.Command {
  return []cli.Command{
    {
      Name:        "team-client-add",
      Usage:       "team-client-add --team <team> --client_id <client_id> --name <name> --redirect_uri <redirect_uri>\n   Your new client has been created.",
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
        defer RecoverStrategy("team-client-add")()
        client := &Client{
          Id   :       c.String("client_id"),
          Name      :  c.String("name"),
          RedirectUri :c.String("redirect_uri"),
          Team   :     c.String("team"),
          client:          NewHTTPClient(&http.Client{}),
        }
        result := client.save()
        fmt.Println(result)
      },
    },
    {
      Name:        "team-client-remove",
      Usage:       "team-client-remove --client_id <client_id> --team <team>\n   The client `<name>` has been deleted.",
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
          defer RecoverStrategy("team-client-remove")()
          client := &Client{
            Id: c.String("client_id"),
            Team: c.String("team"),
            client:    NewHTTPClient(&http.Client{}),
          }
          result := client.remove()
          fmt.Println(result)
        }
      },
    },
  }
}

func (c *Client) save() string {
  path := fmt.Sprintf("/api/teams/%s/clients", c.Team)
  client := &Client{}
  response, err := c.client.MakePost(path, c, client)
  if err != nil {
    return err.Error()
  }

  if response.StatusCode == http.StatusCreated {
    return "Your new client has been created."
  }
  panic("")
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
  return ErrBadRequest.Error()
}
