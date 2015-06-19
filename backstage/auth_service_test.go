package backstage_test

import (
	"github.com/backstage/backstage-cli/backstage"
	"github.com/tsuru/tsuru/fs/fstest"
	. "gopkg.in/check.v1"
)

func (s *S) TestLoginWithValidCredentials(c *C) {
	u, err := userService.Create("Alice", "alice", "alice@example.org", "123")
	t, err := authService.Login(u.Email, u.Password)

	c.Check(err, IsNil)
	c.Assert(t.Type, Equals, "Token")
	c.Assert(t.Expires, Equals, 86400)
	c.Assert(t.CreatedAt, Equals, "2015-05-29T01:05:45Z")
	c.Assert(t.Token, Equals, "RpOMQwiTMtxH6abgwonjBrVhBlrE1jbOxsk86UD_trI=")
}

func (s *S) TestLoginWithInvalidCredentials(c *C) {
	_, err := authService.Login("invalid-email", "invalid-password")
	e := err.(backstage.ResponseError)
	c.Check(err, Not(IsNil))
	c.Assert(e.Error(), Equals, "Authentication failed.")
}

func (s *S) TestLogout(c *C) {
	rfs := &fstest.RecordingFs{}
	backstage.Fsystem = rfs
	defer func() {
		backstage.Fsystem = nil
	}()

	err := authService.Logout()
	c.Check(err, IsNil)
	c.Assert(rfs.HasAction("remove "+backstage.TokenFileName), Equals, true)
}

func (s *S) TestChangePassword(c *C) {
	u, err := userService.Create("Alice", "alice", "alice@example.org", "123")
	c.Check(err, IsNil)
	err = authService.ChangePassword(u.Email, u.Password, "abc", "abc")
	c.Check(err, IsNil)
}

func (s *S) TestChangePasswordWithInvalidConfirmation(c *C) {
	u, err := userService.Create("Alice", "alice", "alice@example.org", "123")
	c.Check(err, IsNil)
	err = authService.ChangePassword(u.Email, u.Password, "abc", "def")
	e := err.(backstage.ResponseError)
	c.Check(err, Not(IsNil))
	c.Assert(e.Error(), Equals, "Your new password and confirmation password do not match.")
}
