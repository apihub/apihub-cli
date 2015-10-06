package apihub_test

import (
	"github.com/apihub/apihub-cli/apihub"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateService(c *C) {
	team, err := teamService.Create("ApiHub Team", "apihub")
	service, err := serviceService.Create("subdomain", false, "description", "documentation", "http://example.org", team.Alias, 10, []string{"XMLTOJSON"})

	c.Check(err, IsNil)
	c.Assert(service.Subdomain, Equals, "subdomain")
	c.Assert(service.Disabled, Equals, false)
	c.Assert(service.Description, Equals, "description")
	c.Assert(service.Documentation, Equals, "documentation")
	c.Assert(service.Endpoint, Equals, "http://example.org")
	c.Assert(service.Team, Equals, team.Alias)
	c.Assert(service.Timeout, Equals, int64(10))
	c.Assert(service.Transformers[0], Equals, "XMLTOJSON")
}

func (s *S) TestCreateServiceMissingRequiredFields(c *C) {
	_, err := serviceService.Create("", false, "", "", "", "team", 10, []string{""})
	e := err.(apihub.ResponseError)
	c.Assert(e.Error(), Equals, "Team not found.")
}

func (s *S) TestUpdateService(c *C) {
	team, err := teamService.Create("ApiHub Team", "apihub")
	service, err := serviceService.Create("subdomain", false, "description", "documentation", "http://example.org", team.Alias, 10, []string{"XMLTOJSON"})

	service, err = serviceService.Update("subdomain", true, "new description", "new documentation", "http://example.org/v2", team.Alias, 1, []string{"XMLTOJSON"})
	c.Check(err, IsNil)
	c.Assert(service.Description, Equals, "new description")
}

func (s *S) TestUpdateServiceNotFound(c *C) {
	team, err := teamService.Create("ApiHub Team", "apihub")
	_, err = serviceService.Update("not-found", true, "new description", "new documentation", "http://example.org/v2", team.Alias, 1, []string{"XMLTOJSON"})
	e := err.(apihub.ResponseError)
	c.Assert(e.Error(), Equals, "Service not found.")
}

func (s *S) TestDeleteService(c *C) {
	team, err := teamService.Create("ApiHub Team", "apihub")
	service, err := serviceService.Create("subdomain", false, "description", "documentation", "http://example.org", team.Alias, 10, []string{"XMLTOJSON"})

	err = serviceService.Delete(service.Subdomain, service.Team)
	c.Check(err, IsNil)
}
