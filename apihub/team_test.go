package apihub_test

import (
	"github.com/apihub/apihub-cli/apihub"
	. "gopkg.in/check.v1"
)

func (s *S) TestContainsUserByEmail(c *C) {
	team := apihub.Team{
		Users: []string{"alice@example.org"},
	}

	i, ok := team.ContainsUserByEmail("alice@example.org")
	c.Assert(ok, Equals, true)
	c.Assert(i, Equals, 0)
}

func (s *S) TestContainsUserByEmailWithNotFound(c *C) {
	team := apihub.Team{
		Users: []string{},
	}

	i, ok := team.ContainsUserByEmail("alice@example.org")
	c.Assert(ok, Equals, false)
	c.Assert(i, Equals, -1)
}
