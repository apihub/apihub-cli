package backstage_test

import (
	"github.com/backstage/backstage-cli/maestro"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateUser(c *C) {
	u, err := userService.Create("Alice", "alice", "alice@example.org", "123")

	c.Check(err, IsNil)
	c.Assert(u.Name, Equals, "Alice")
	c.Assert(u.Username, Equals, "alice")
	c.Assert(u.Email, Equals, "alice@example.org")
	c.Assert(u.Password, Equals, "123")
}

func (s *S) TestCreateUserMissingRequiredFields(c *C) {
	_, err := userService.Create("", "", "", "")
	e := err.(backstage.ResponseError)
	c.Assert(e.Error(), Equals, "Name/Email/Username/Password cannot be empty.")
}

func (s *S) TestDeleteUser(c *C) {
	u, err := userService.Create("Alice", "alice", "alice@example.org", "123")
	_, err = authService.Login(u.Email, u.Password)

	err = userService.Delete()
	c.Check(err, IsNil)
}
