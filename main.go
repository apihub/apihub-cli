package main

import (
	"os"

	"github.com/backstage/backstage-cli/commands"
	"github.com/backstage/backstage-cli/maestro"
	"github.com/codegangsta/cli"
)

func main() {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate
	app := cli.NewApp()
	app.Name = "backstage"
	app.Usage = "An open source solution for publishing APIs."
	app.Version = backstage.BackstageClientVersion
	app.HideHelp = true
	app.EnableBashCompletion = true

	currentTarget, err := backstage.GetCurrentTarget()
	if err != nil {
		panic("Your target file is corrupted. Please delete it and add your target. Sorry about that.")
	}

	httpClient := backstage.NewHTTPClient(currentTarget)

	m := NewManager(app)
	m.Register(&commands.Auth{Service: backstage.NewAuthService(httpClient)})
	m.Register(&commands.App{Service: backstage.NewAppService(httpClient)})
	m.Register(&commands.Service{Service: backstage.NewServiceService(httpClient)})
	m.Register(&commands.Target{})
	m.Register(&commands.Team{Service: backstage.NewTeamService(httpClient)})
	m.Register(&commands.User{Service: backstage.NewUserService(httpClient)})

	app.Run(os.Args)
}
