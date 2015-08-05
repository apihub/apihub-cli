package apihub_test

import (
	"github.com/apihub/apihub-cli/maestro"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateTeam(c *C) {
	t, err := teamService.Create("ApiHub Team", "apihub")

	c.Check(err, IsNil)
	c.Assert(t.Name, Equals, "ApiHub Team")
	c.Assert(t.Alias, Equals, "apihub")
	c.Assert(t.Owner, Equals, "alice@example.org")
}

func (s *S) TestCreateTeamMissingRequiredFields(c *C) {
	_, err := teamService.Create("", "apihub")
	e := err.(apihub.ResponseError)
	c.Assert(e.Error(), Equals, "Name cannot be empty.")
}

func (s *S) TestUpdateTeam(c *C) {
	_, err := teamService.Create("ApiHub Team", "apihub")

	t, err := teamService.Update("New Name", "apihub")

	c.Check(err, IsNil)
	c.Assert(t.Name, Equals, "New Name")
	c.Assert(t.Alias, Equals, "apihub")
}

func (s *S) TestUpdateTeamWithInvalidAlais(c *C) {
	_, err := teamService.Update("New Name", "invalid-alias")
	e := err.(apihub.ResponseError)
	c.Assert(e.Error(), Equals, "Team not found.")
}

func (s *S) TestUpdateTeamWhithoutBeingPartOfTeam(c *C) {
	_, err := teamService.Create("ApiHub Team", "apihub")

	t, err := teamService.Update("New Name", "apihub")

	c.Check(err, IsNil)
	c.Assert(t.Name, Equals, "New Name")
	c.Assert(t.Alias, Equals, "apihub")
}

func (s *S) TestTeamInfo(c *C) {
	_, err := teamService.Create("ApiHub Team", "apihub")

	t, err := teamService.Info("apihub")
	c.Check(err, IsNil)
	c.Assert(t.Name, Equals, "ApiHub Team")
	c.Assert(t.Alias, Equals, "apihub")
	c.Assert(t.Owner, Equals, "alice@example.org")
}

func (s *S) TestTeamInfoNotFound(c *C) {
	_, err := teamService.Info("not-found")
	e := err.(apihub.ResponseError)
	c.Assert(e.Error(), Equals, "Team not found.")
}

func (s *S) TestListTeam(c *C) {
	teamService.Create("ApiHub Team", "apihub")
	teamService.Create("Varnish", "varnish")
	teamService.Create("Oxi", "oxi")

	teams, err := teamService.List()
	c.Check(err, IsNil)
	c.Assert(len(teams), Equals, 3)
}

func (s *S) TestDeleteTeam(c *C) {
	teamService.Create("ApiHub Team", "apihub")
	t, _ := teamService.Info("apihub")
	c.Assert(t.Name, Equals, "ApiHub Team")

	err := teamService.Delete("apihub")
	c.Check(err, IsNil)
	t, _ = teamService.Info("apihub")
	c.Assert(t.Name, Equals, "")
}

func (s *S) TestTeamAddUser(c *C) {
	t, _ := teamService.Create("ApiHub Team", "apihub")
	ok, err := teamService.AddUser(t.Alias, "bob@example.org")
	team, _ := teamService.Info(t.Alias)
	c.Assert(len(team.Users), Equals, 2)
	c.Check(err, IsNil)
	c.Assert(ok, Equals, true)
}

func (s *S) TestTeamAddUserWithNotFoundTeam(c *C) {
	ok, err := teamService.AddUser("not-found", "alice@example.org")
	c.Check(err, Not(IsNil))
	c.Assert(ok, Equals, false)
}

func (s *S) TestTeamRemoveUser(c *C) {
	t, _ := teamService.Create("ApiHub Team", "apihub")
	teamService.AddUser(t.Alias, "bob@example.org")
	team, _ := teamService.Info(t.Alias)
	c.Assert(len(team.Users), Equals, 2)

	ok, err := teamService.RemoveUser(t.Alias, "bob@example.org")
	c.Check(err, IsNil)
	c.Assert(ok, Equals, true)
	team, _ = teamService.Info(t.Alias)
	c.Assert(len(team.Users), Equals, 1)
}

func (s *S) TestTeamRemoveUserWithNotFoundTeam(c *C) {
	ok, err := teamService.RemoveUser("not-found", "alice@example.org")
	c.Check(err, Not(IsNil))
	c.Assert(ok, Equals, false)
}
