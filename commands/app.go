package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/backstage/backstage-cli/maestro"
	"github.com/codegangsta/cli"
)

type App struct {
	Service *backstage.AppService
}

func (cmd *App) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "app-create",
			Usage:       "app-create --team <team> --client_id <client_id> --client_secret <client_secret> --name <name> --redirect_uris <redirect_uris>",
			Description: "Create a new app.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "client_id, i", Value: "", Usage: "App id (used by oAuth 2.0)"},
				cli.StringFlag{Name: "client_secret, s", Value: "", Usage: "Secret Key (used by oAuth 2.0)"},
				cli.StringFlag{Name: "name, n", Value: "", Usage: "App name"},
				cli.StringFlag{Name: "redirect_uris, r", Value: "", Usage: "App Redirect Uris, comma separated (used by oAuth 2.0)"},
				cli.StringFlag{Name: "team, t", Value: "", Usage: "Team"},
			},
			Action: cmd.appCreate,
		},

		{
			Name:        "app-update",
			Usage:       "app-update --team <team> --client_id <client_id> --client_secret <client_secret>  --name <name> --redirect_uris <redirect_uris>",
			Description: "Update an existing app.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "client_id, i", Value: "", Usage: "App id (used by oAuth 2.0)"},
				cli.StringFlag{Name: "client_secret, s", Value: "", Usage: "Secret Key (used by oAuth 2.0)"},
				cli.StringFlag{Name: "name, n", Value: "", Usage: "App name"},
				cli.StringFlag{Name: "redirect_uris, r", Value: "", Usage: "App Redirect Uri (used by oAuth 2.0)"},
				cli.StringFlag{Name: "team, t", Value: "", Usage: "Team"},
			},
			Action: cmd.appUpdate,
		},

		{
			Name:        "app-remove",
			Usage:       "app-remove --client_id <client_id> --team <team>\n   The app `<name>` has been deleted.",
			Description: "Remove an existing app.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "client_id, i", Value: "", Usage: "App Id"},
				cli.StringFlag{Name: "team, t", Value: "", Usage: "Team"},
			},
			Action: cmd.appRemove,
		},

		{
			Name:        "app-info",
			Usage:       "app-info --client_id <client_id>",
			Description: "Retrieve app information.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "client_id, i", Value: "", Usage: "App Id"},
			},
			Action: cmd.appInfo,
		},
	}
}

func (cmd *App) appCreate(c *cli.Context) {
	defer RecoverStrategy("app-create")()
	uris := strings.Split(c.String("redirect_uris"), ",")
	_, err := cmd.Service.Create(c.String("team"), c.String("client_id"), c.String("name"), uris, c.String("client_secret"))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Your new app has been created.")
	}
}

func (cmd *App) appUpdate(c *cli.Context) {
	defer RecoverStrategy("app-update")()
	uris := strings.Split(c.String("redirect_uris"), ",")
	_, err := cmd.Service.Update(c.String("team"), c.String("client_id"), c.String("name"), uris, c.String("client_secret"))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Your new app has been updated.")
	}
}

func (cmd *App) appRemove(c *cli.Context) {
	defer RecoverStrategy("app-remove")()

	context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
	if Confirm(context, "Are you sure you want to delete this app? This action cannot be undone.") != true {
		fmt.Println(backstage.ErrCommandCancelled)
	} else {
		err := cmd.Service.Delete(c.String("team"), c.String("client_id"))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("The app `" + c.String("client_id") + "` has been deleted.")
		}
	}
}

func (cmd *App) appInfo(c *cli.Context) {
	defer RecoverStrategy("app-info")()

	appID := c.String("client_id")
	if appID == "" {
		appID = c.Args().First()
	}
	app, err := cmd.Service.Info(appID)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Team Name: " + app.Team)
		fmt.Println("")
		appTable := &Table{
			Title:   "Apps Details:",
			Content: [][]string{},
			Header:  []string{"Name", "Redirect Uri", "Id", "Secret"},
		}

		line := []string{}
		line = append(line, app.Name)
		line = append(line, strings.Join(app.RedirectURIs, ", "))
		line = append(line, app.ClientID)
		line = append(line, app.ClientSecret)
		appTable.Content = append(appTable.Content, line)

		context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
		appTable.Render(context)
		fmt.Println()
	}
}
