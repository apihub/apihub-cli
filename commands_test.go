package main

import (
	"github.com/codegangsta/cli"
	. "gopkg.in/check.v1"
)

type TestCommand struct{}

func (t *TestCommand) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "login",
			Action: func(c *cli.Context) {},
		},
	}
}

func (s *S) TestRegisterCommands(c *C) {
	app := cli.NewApp()
	c.Assert(len(app.Commands), Equals, 0)
	m := NewManager(app)
	m.Register(&TestCommand{})
	c.Assert(len(app.Commands), Equals, 1)
}
