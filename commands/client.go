package commands

import (
	"fmt"
	"os"

	"github.com/backstage/backstage-client/backstage"
	"github.com/codegangsta/cli"
)

type Client struct {
	Service *backstage.ClientService
}

func (cmd *Client) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "client-create",
			Usage:       "client-create --team <team> --client_id <client_id> --secret <secret> --name <name> --redirect_uri <redirect_uri>",
			Description: "Create a new client.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "client_id, i", Value: "", Usage: "Client id (used by oAuth 2.0)"},
				cli.StringFlag{Name: "secret, s", Value: "", Usage: "Secret Key (used by oAuth 2.0)"},
				cli.StringFlag{Name: "name, n", Value: "", Usage: "Client name"},
				cli.StringFlag{Name: "redirect_uri, r", Value: "", Usage: "Client Redirect Uri (used by oAuth 2.0)"},
				cli.StringFlag{Name: "team, t", Value: "", Usage: "Team"},
			},
			Action: cmd.clientCreate,
		},

		{
			Name:        "client-update",
			Usage:       "client-update --team <team> --client_id <client_id> --secret <secret>  --name <name> --redirect_uri <redirect_uri>",
			Description: "Update an existing client.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "client_id, i", Value: "", Usage: "Client id (used by oAuth 2.0)"},
				cli.StringFlag{Name: "secret, s", Value: "", Usage: "Secret Key (used by oAuth 2.0)"},
				cli.StringFlag{Name: "name, n", Value: "", Usage: "Client name"},
				cli.StringFlag{Name: "redirect_uri, r", Value: "", Usage: "Client Redirect Uri (used by oAuth 2.0)"},
				cli.StringFlag{Name: "team, t", Value: "", Usage: "Team"},
			},
			Action: cmd.clientUpdate,
		},

		{
			Name:        "client-remove",
			Usage:       "client-remove --client_id <client_id> --team <team>\n   The client `<name>` has been deleted.",
			Description: "Remove an existing client.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "client_id, i", Value: "", Usage: "Client Id"},
				cli.StringFlag{Name: "team, t", Value: "", Usage: "Team"},
			},
			Action: cmd.clientRemove,
		},

		{
			Name:        "client-info",
			Usage:       "client-info --client_id <client_id>",
			Description: "Retrieve client information.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "client_id, i", Value: "", Usage: "Client Id"},
			},
			Action: cmd.clientInfo,
		},
	}
}

func (cmd *Client) clientCreate(c *cli.Context) {
	defer RecoverStrategy("client-create")()

	_, err := cmd.Service.Create(c.String("team"), c.String("client_id"), c.String("name"), c.String("redirect_uri"), c.String("secret"))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Your new client has been created.")
	}
}

func (cmd *Client) clientUpdate(c *cli.Context) {
	defer RecoverStrategy("client-update")()

	_, err := cmd.Service.Update(c.String("team"), c.String("client_id"), c.String("name"), c.String("redirect_uri"), c.String("secret"))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Your new client has been updated.")
	}
}

func (cmd *Client) clientRemove(c *cli.Context) {
	defer RecoverStrategy("client-remove")()

	context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
	if Confirm(context, "Are you sure you want to delete this client? This action cannot be undone.") != true {
		fmt.Println(backstage.ErrCommandCancelled)
	} else {
		err := cmd.Service.Delete(c.String("team"), c.String("client_id"))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("The client `" + c.String("client_id") + "` has been deleted.")
		}
	}
}

func (cmd *Client) clientInfo(c *cli.Context) {
	defer RecoverStrategy("client-info")()

	client_id := c.String("client_id")
	if client_id == "" {
		client_id = c.Args().First()
	}
	client, err := cmd.Service.Info(client_id)
	if err != nil {
		fmt.Println(err.Error())
	} else {
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

		context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
		clientTable.Render(context)
		fmt.Println("\n")
	}
}
