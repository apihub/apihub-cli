package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/backstage/backstage-cli/maestro"
	"github.com/codegangsta/cli"
)

type Service struct {
	Service *backstage.ServiceService
}

func (cmd *Service) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "service-create",
			Usage:       "service-create --team <team> --subdomain <subdomain> --endpoint <api_endpoint> --timeout <timeout> --transformers <transformers>",
			Description: "Create a new service.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "description, desc", Value: "", Usage: "Service description"},
				cli.StringFlag{Name: "disabled, dis", Value: "", Usage: "Disable the service"},
				cli.StringFlag{Name: "documentation, doc", Value: "", Usage: "Url with the documentation"},
				cli.StringFlag{Name: "endpoint, e", Value: "", Usage: "Url where the service is running"},
				cli.StringFlag{Name: "subdomain, s", Value: "", Usage: "Desired subdomain"},
				cli.StringFlag{Name: "team, t", Value: "", Usage: "Team responsible for the service"},
				cli.StringFlag{Name: "timeout", Value: "", Usage: "Timeout"},
				cli.StringFlag{Name: "transformers, tf", Value: "", Usage: "Transformers"},
			},
			Action: cmd.userCreate,
		},

		{
			Name:        "service-update",
			Usage:       "service-update --team <team> --subdomain <subdomain> --endpoint <api_endpoint> --timeout <timeout> --transformers <transformers>",
			Description: "Update an existing service.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "description, desc", Value: "", Usage: "Service description"},
				cli.StringFlag{Name: "disabled, dis", Value: "", Usage: "Disable the service"},
				cli.StringFlag{Name: "documentation, doc", Value: "", Usage: "Url with the documentation"},
				cli.StringFlag{Name: "endpoint, e", Value: "", Usage: "Url where the service is running"},
				cli.StringFlag{Name: "subdomain, s", Value: "", Usage: "Desired subdomain"},
				cli.StringFlag{Name: "team, t", Value: "", Usage: "Team responsible for the service"},
				cli.StringFlag{Name: "timeout", Value: "", Usage: "Timeout"},
				cli.StringFlag{Name: "transformers, tf", Value: "", Usage: "Transformers"},
			},
			Action: cmd.userUpdate,
		},

		{
			Name:        "service-remove",
			Usage:       "service-remove --subdomain <subdomain>",
			Description: "Remove an existing service.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "subdomain, s", Value: "", Usage: "Subdomain"},
				cli.StringFlag{Name: "team, t", Value: "", Usage: "Team"},
			},
			Action: cmd.serviceRemove,
		},
	}
}

func (cmd *Service) userCreate(c *cli.Context) {
	defer RecoverStrategy("service-create")()

	subdomain, disabled, description, documentation, endpoint, team, timeout, transformers := cmd.parseArgs(c)
	_, err := cmd.Service.Create(subdomain, disabled, description, documentation, endpoint, team, timeout, transformers)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Your new service has been created.")
	}
}

func (cmd *Service) userUpdate(c *cli.Context) {
	defer RecoverStrategy("service-update")()

	subdomain, disabled, description, documentation, endpoint, team, timeout, transformers := cmd.parseArgs(c)
	_, err := cmd.Service.Update(subdomain, disabled, description, documentation, endpoint, team, timeout, transformers)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Your new service has been updated.")
	}
}

func (cmd *Service) serviceRemove(c *cli.Context) {
	defer RecoverStrategy("service-remove")()

	context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
	if Confirm(context, "Are you sure you want to delete this service? This action cannot be undone.") != true {
		fmt.Println(backstage.ErrCommandCancelled)
	} else {
		err := cmd.Service.Delete(c.String("subdomain"), c.String("team"))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("The service `" + c.String("subdomain") + "` has been deleted.")
		}
	}
}

func (cmd *Service) parseArgs(c *cli.Context) (subdomain string, disabled bool, description string, documentation string, endpoint string, team string, timeout int64, transformers []string) {
	disabled, err := strconv.ParseBool(c.String("disabled"))
	if err != nil {
		disabled = false
	}
	timeout, err = strconv.ParseInt(c.String("timeout"), 10, 0)
	if err != nil {
		timeout = 0
	}
	var transf []string
	if c.String("transformers") != "" {
		transf = strings.Split(c.String("transformers"), ",")
	}
	return c.String("subdomain"), disabled, c.String("description"), c.String("documentation"), c.String("endpoint"), c.String("team"), int64(timeout), transf
}
