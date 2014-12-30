package main

import (
	"net/http"

	ttesting "github.com/tsuru/tsuru/cmd/testing"
	"github.com/tsuru/tsuru/fs/testing"
	. "gopkg.in/check.v1"
)

func (s *S) TestUserCreate(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	user := &User{
		Name:     "Alice",
		Email:    "alice@example.org",
		Username: "alice",
		Password: "123",
	}
	transport := ttesting.Transport{
		Status:  http.StatusCreated,
		Message: `{"name":"` + user.Name + `","email":"` + user.Email + `","username":"` + user.Username + `"}`,
	}
	user.client = NewHTTPClient(&http.Client{Transport: &transport})
	r := user.save()
	c.Assert(r, Equals, "Your account has been created.")
}

func (s *S) TestUserCreateInvalidUserInfo(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	transport := ttesting.Transport{
		Status:  http.StatusBadRequest,
		Message: `{"error":"bad_request","error_description":"Someone already has that username. Could you try another?"}`,
	}
	user := &User{
		Name:     "Alice",
		Email:    "alice@example.org",
		Username: "alice",
		Password: "123",
		client:   NewHTTPClient(&http.Client{Transport: &transport}),
	}
	r := user.save()
	c.Assert(r, Equals, "Someone already has that username. Could you try another?")
}

func (s *S) TestUserRemove(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	user := &User{
		Name:     "Alice",
		Email:    "alice@example.org",
		Username: "alice",
		Password: "123",
	}
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: `{"name":"` + user.Name + `","email":"` + user.Email + `","username":"` + user.Username + `"}`,
	}
	user.client = NewHTTPClient(&http.Client{Transport: &transport})
	r := user.remove()
	c.Assert(r, Equals, "Your account has been deleted.")
}

func (s *S) TestUserRemoveWithoutTarget(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current:\n"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	user := &User{
		Name:     "Alice",
		Email:    "alice@example.org",
		Username: "alice",
		Password: "123",
	}
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: `{}`,
	}
	user.client = NewHTTPClient(&http.Client{Transport: &transport})
	r := user.remove()
	c.Assert(r, Equals, "You have not selected any target as default. For more details, please run `backstage target-set -h`.")
}