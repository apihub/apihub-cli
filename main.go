package main

import (
	"os"

	"github.com/apihub/apihub-cli/apihub"
	"github.com/apihub/apihub-cli/commands"
	"github.com/codegangsta/cli"
)

func main() {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate
	app := cli.NewApp()
	app.Name = "apihub"
	app.Usage = "An open source solution for publishing APIs."
	app.Version = apihub.ApiHubClientVersion
	app.HideHelp = true
	app.EnableBashCompletion = true

	currentTarget, _ := apihub.GetCurrentTarget()
	httpClient := apihub.NewHTTPClient(currentTarget)

	m := NewManager(app)
	m.Register(&commands.Auth{Service: apihub.NewAuthService(httpClient)})
	m.Register(&commands.App{Service: apihub.NewAppService(httpClient)})
	m.Register(&commands.Service{Service: apihub.NewServiceService(httpClient)})
	m.Register(&commands.Target{})
	m.Register(&commands.Team{Service: apihub.NewTeamService(httpClient)})
	m.Register(&commands.User{Service: apihub.NewUserService(httpClient)})

	app.Run(os.Args)
}
