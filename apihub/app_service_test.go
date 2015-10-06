package apihub_test

import (
	"github.com/apihub/apihub-cli/apihub"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateApp(c *C) {
	app, err := appService.Create("apihub", "123", "Documents", []string{"http://example.org/auth"}, "super-secret")

	c.Check(err, IsNil)
	c.Assert(app.Team, Equals, "apihub")
	c.Assert(app.ClientID, Equals, "123")
	c.Assert(app.Name, Equals, "Documents")
	c.Assert(app.RedirectURIs, DeepEquals, []string{"http://example.org/auth"})
	c.Assert(app.ClientSecret, Equals, "super-secret")
}

func (s *S) TestCreateAppMissingRequiredFields(c *C) {
	_, err := appService.Create("apihub", "123", "", []string{"http://example.org/auth"}, "super-secret")
	e := err.(apihub.ResponseError)
	c.Assert(e.Error(), Equals, "Name cannot be empty.")
}

func (s *S) TestUpdateApp(c *C) {
	t, err := teamService.Create("ApiHub Team", "apihub")
	app, _ := appService.Create(t.Alias, "123", "Documents", []string{"http://example.org/auth"}, "super-secret")

	_, err = appService.Update(t.Alias, app.ClientID, "New Name", []string{"http://example.org/v2/auth"}, "amazing")
	c.Check(err, IsNil)

	app, _ = appService.Info("123")
	c.Assert(app.Team, Equals, t.Alias)
	c.Assert(app.ClientID, Equals, "123")
	c.Assert(app.Name, Equals, "New Name")
	c.Assert(app.RedirectURIs, DeepEquals, []string{"http://example.org/v2/auth"})
	c.Assert(app.ClientSecret, Equals, "amazing")
}

func (s *S) TestUpdateAppNotFound(c *C) {
	t, err := teamService.Create("ApiHub Team", "apihub")
	_, err = appService.Update(t.Alias, "123", "New Name", []string{"http://example.org/v2/auth"}, "amazing")
	e := err.(apihub.ResponseError)
	c.Assert(e.Error(), Equals, "App not found.")
}

func (s *S) TestAppInfo(c *C) {
	_, err := appService.Create("apihub", "123", "Documents", []string{"http://example.org/auth"}, "super-secret")
	c.Check(err, IsNil)

	app, err := appService.Info("123")
	c.Check(err, IsNil)
	c.Assert(app.Team, Equals, "apihub")
	c.Assert(app.ClientID, Equals, "123")
	c.Assert(app.RedirectURIs, DeepEquals, []string{"http://example.org/auth"})
	c.Assert(app.ClientSecret, Equals, "super-secret")
}

func (s *S) TestAppInfoNotFound(c *C) {
	_, err := appService.Info("not-found")
	e := err.(apihub.ResponseError)
	c.Assert(e.Error(), Equals, "App not found.")
}

func (s *S) TestDeleteApp(c *C) {
	app, err := appService.Create("apihub", "123", "Documents", []string{"http://example.org/auth"}, "super-secret")
	found, err := appService.Info(app.ClientID)
	c.Assert(found.Name, Equals, app.Name)

	err = appService.Delete(app.Team, app.ClientID)
	c.Check(err, IsNil)
	app, _ = appService.Info(app.ClientID)
	c.Assert(app.Name, Equals, "")
}
