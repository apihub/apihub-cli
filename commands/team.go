package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/backstage/backstage-cli/maestro"
	"github.com/codegangsta/cli"
)

type Team struct {
	Service *backstage.TeamService
}

func (cmd *Team) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "team-create",
			Usage:       "team-create --name <name>",
			Description: "Create a team.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "name, n", Value: "", Usage: "Name of the team"},
				cli.StringFlag{Name: "alias, a", Value: "", Usage: "Alias"},
			},
			Action: cmd.teamCreate,
		},

		{
			Name:        "team-update",
			Usage:       "team-update --name <name>",
			Description: "Update an existing.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "name, n", Value: "", Usage: "Name of the team"},
				cli.StringFlag{Name: "alias, a", Value: "", Usage: "Alias"},
			},
			Action: cmd.teamUpdate,
		},

		{
			Name:        "team-info",
			Usage:       "team-info --alias <alias>",
			Description: "Return team info and lists of members and services.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "alias, a", Value: "", Usage: "Team alias"},
			},
			Action: cmd.teamInfo,
		},

		{
			Name:        "team-list",
			Usage:       "team-list",
			Description: "Return a list of all teams.",
			Action:      cmd.teamList,
		},

		{
			Name:        "team-remove",
			Usage:       "team-remove --alias <alias>",
			Description: "Delete a team.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "alias, a", Value: "", Usage: "Team alias"},
			},
			Action: cmd.teamRemove,
		},

		{
			Name:        "team-user-add",
			Usage:       "team-user-add --team <team-alias> --email <email>",
			Description: "Add a user to a team.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "team, t", Value: "", Usage: "Name of the team"},
				cli.StringFlag{Name: "email, e", Value: "", Usage: "User's email"},
			},
			Action: cmd.teamUserAdd,
		},

		{
			Name:        "team-user-remove",
			Usage:       "team-user-remove --team <team-alias> --email <email>",
			Description: "Remove a user from a team.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "team, t", Value: "", Usage: "Name of the team"},
				cli.StringFlag{Name: "email, e", Value: "", Usage: "User's email"},
			},
			Action: cmd.teamUserRemove,
		},
	}
}

func (cmd *Team) teamCreate(c *cli.Context) {
	defer RecoverStrategy("team-create")()

	name := cmd.name(c)
	_, err := cmd.Service.Create(name, c.String("alias"))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Your new team has been created.")
	}
}

func (cmd *Team) teamUpdate(c *cli.Context) {
	defer RecoverStrategy("team-update")()

	name := cmd.name(c)
	_, err := cmd.Service.Update(name, c.String("alias"))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Your team has been updated.")
	}
}

func (cmd *Team) teamRemove(c *cli.Context) {
	defer RecoverStrategy("team-remove")()

	context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
	if Confirm(context, "Are you sure you want to delete this team? This action cannot be undone.") != true {
		fmt.Println(backstage.ErrCommandCancelled)
	} else {
		alias := c.String("alias")
		if alias == "" {
			alias = c.Args().First()
		}

		err := cmd.Service.Delete(alias)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Your team has been deleted.")
		}
	}
}

func (cmd *Team) teamUserAdd(c *cli.Context) {
	defer RecoverStrategy("team-user-add")()

	_, err := cmd.Service.AddUser(c.String("team"), c.String("email"))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("User `" + c.String("email") + "` has been added successfully to team `" + c.String("team") + "`.")
	}
}

func (cmd *Team) teamUserRemove(c *cli.Context) {
	defer RecoverStrategy("team-user-remove")()

	_, err := cmd.Service.RemoveUser(c.String("team"), c.String("email"))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("User `" + c.String("email") + "` has been removed successfully from the team `" + c.String("team") + "`.")
	}
}

func (cmd *Team) teamInfo(c *cli.Context) {
	defer RecoverStrategy("team-info")()

	alias := cmd.alias(c)
	team, err := cmd.Service.Info(alias)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Name: " + team.Name)
		fmt.Println("Alias: " + team.Alias)
		fmt.Println("Owner: " + team.Owner + "\n")

		var tables []*Table

		// Team Members table
		membersTable := &Table{
			Content: [][]string{},
			Header:  []string{"Team Members"},
		}
		for _, member := range team.Users {
			line := []string{}
			line = append(line, member)
			membersTable.Content = append(membersTable.Content, line)
		}
		tables = append(tables, membersTable)

		//Services table
		if len(team.Services) > 0 {
			servicesTable := &Table{
				Title:   "Available Services:",
				Content: [][]string{},
				Header:  []string{"Subdomain", "Endpoint", "Owner"},
			}
			for _, service := range team.Services {
				line := []string{}
				line = append(line, service.Subdomain)
				line = append(line, service.Endpoint)
				line = append(line, service.Owner)
				servicesTable.Content = append(servicesTable.Content, line)
			}
			tables = append(tables, servicesTable)
		}

		// Apps table
		if len(team.Apps) > 0 {
			appsTable := &Table{
				Title:   "Available Apps:",
				Content: [][]string{},
				Header:  []string{"Id", "Name", "Redirect Uri"},
			}
			for _, app := range team.Apps {
				line := []string{}
				line = append(line, app.ClientID)
				line = append(line, app.Name)
				line = append(line, strings.Join(app.RedirectURIs, ", "))
				appsTable.Content = append(appsTable.Content, line)
			}
			tables = append(tables, appsTable)
		}

		var context *Context
		for _, table := range tables {
			context = &Context{Stdout: os.Stdout, Stdin: os.Stdin}
			table.Render(context)
			fmt.Println()
		}
	}
}

func (cmd *Team) teamList(c *cli.Context) {
	defer RecoverStrategy("team-list")()

	list, err := cmd.Service.List()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if len(list) > 0 {
			table := &Table{
				Content: [][]string{},
				Header:  []string{"Team Name", "Alias", "Owner"},
			}
			for _, team := range list {
				line := []string{}
				line = append(line, team.Name, team.Alias, team.Owner)
				table.Content = append(table.Content, line)
			}
			context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
			table.Render(context)
		} else {
			fmt.Println("You have no teams.")
		}
	}
}

func (cmd *Team) name(c *cli.Context) string {
	name := c.String("name")
	if name == "" {
		name = c.Args().First()
	}
	return name
}

func (cmd *Team) alias(c *cli.Context) string {
	alias := c.String("alias")
	if alias == "" {
		alias = c.Args().First()
	}
	return alias
}
