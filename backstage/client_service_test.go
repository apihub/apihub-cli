package backstage_test

import (
	"github.com/backstage/backstage-client/backstage"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateClient(c *C) {
	cli, err := clientService.Create("backstage", "123", "Documents", "http://example.org/auth", "super-secret")

	c.Check(err, IsNil)
	c.Assert(cli.Team, Equals, "backstage")
	c.Assert(cli.Id, Equals, "123")
	c.Assert(cli.Name, Equals, "Documents")
	c.Assert(cli.RedirectUri, Equals, "http://example.org/auth")
	c.Assert(cli.Secret, Equals, "super-secret")
}

func (s *S) TestCreateClientMissingRequiredFields(c *C) {
	_, err := clientService.Create("backstage", "123", "", "http://example.org/auth", "super-secret")
	e := err.(backstage.ResponseError)
	c.Assert(e.Error(), Equals, "Name cannot be empty.")
}

func (s *S) TestUpdateClient(c *C) {
	t, err := teamService.Create("Backstage Team", "backstage")
	cli, _ := clientService.Create(t.Alias, "123", "Documents", "http://example.org/auth", "super-secret")

	_, err = clientService.Update(t.Alias, cli.Id, "New Name", "http://example.org/v2/auth", "amazing")
	c.Check(err, IsNil)

	cli, _ = clientService.Info("123")
	c.Assert(cli.Team, Equals, t.Alias)
	c.Assert(cli.Id, Equals, "123")
	c.Assert(cli.Name, Equals, "New Name")
	c.Assert(cli.RedirectUri, Equals, "http://example.org/v2/auth")
	c.Assert(cli.Secret, Equals, "amazing")
}

func (s *S) TestUpdateClientWhenTeamNotFound(c *C) {
	_, err := clientService.Update("backstage", "123", "New Name", "http://example.org/v2/auth", "amazing")
	e := err.(backstage.ResponseError)
	c.Assert(e.Error(), Equals, "Team not found.")
}

func (s *S) TestUpdateClientNotFound(c *C) {
	t, err := teamService.Create("Backstage Team", "backstage")
	_, err = clientService.Update(t.Alias, "123", "New Name", "http://example.org/v2/auth", "amazing")
	e := err.(backstage.ResponseError)
	c.Assert(e.Error(), Equals, "Client not found.")
}

func (s *S) TestClientInfo(c *C) {
	_, err := clientService.Create("backstage", "123", "Documents", "http://example.org/auth", "super-secret")
	c.Check(err, IsNil)

	cli, err := clientService.Info("123")
	c.Check(err, IsNil)
	c.Assert(cli.Team, Equals, "backstage")
	c.Assert(cli.Id, Equals, "123")
	c.Assert(cli.RedirectUri, Equals, "http://example.org/auth")
	c.Assert(cli.Secret, Equals, "super-secret")
}

func (s *S) TestClientInfoNotFound(c *C) {
	_, err := clientService.Info("not-found")
	e := err.(backstage.ResponseError)
	c.Assert(e.Error(), Equals, "Client not found.")
}

func (s *S) TestDeleteClient(c *C) {
	cli, err := clientService.Create("backstage", "123", "Documents", "http://example.org/auth", "super-secret")
	found, err := clientService.Info(cli.Id)
	c.Assert(found.Name, Equals, cli.Name)

	err = clientService.Delete(cli.Team, cli.Id)
	c.Check(err, IsNil)
	cli, _ = clientService.Info(cli.Id)
	c.Assert(cli.Name, Equals, "")
}
