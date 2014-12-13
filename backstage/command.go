package main

import "github.com/codegangsta/cli"

type Command interface {
	GetCommands() []cli.Command
}

type Manager struct {
	App *cli.App
}

func NewManager(app *cli.App) *Manager {
	return &Manager{App: app}
}

func (m *Manager) Register(command Command) {
	cmds := command.GetCommands()
	for _, cmd := range cmds {
		m.App.Commands = append(m.App.Commands, cmd)
	}
}