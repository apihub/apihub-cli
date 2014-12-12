package main

import (
	"net/http"

	ttesting "github.com/tsuru/tsuru/cmd/testing"
	"github.com/tsuru/tsuru/fs/testing"
	. "gopkg.in/check.v1"
)

func (s *S) TestTeamCreate(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	team := &Team{
		Name: "Kotobuki",
	}
	transport := ttesting.Transport{
		Status:  http.StatusCreated,
		Message: `{"id":"548ab5b00904b8bf2e8dd838","name":"Kotobuki","alias":"kotobuki","users":["alice@example.org"],"owner":"alice@example.org"}`,
	}
	team.client = NewClient(&http.Client{Transport: &transport})
	r := team.save()
	c.Assert(r, Equals, "Team created successfully.")
}

func (s *S) TestTeamCreateWithExistingName(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	transport := ttesting.Transport{
		Status:  http.StatusBadRequest,
		Message: `{"status_code":400,"message":"Someone already has that team name. Could you try another?"}`,
	}
	team := &Team{
		Name: "Kotobuki",
	}
	team.client = NewClient(&http.Client{Transport: &transport})
	r := team.save()
	c.Assert(r, Equals, "Someone already has that team name. Could you try another?")
}

func (s *S) TestTeamRemove(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	team := &Team{
		Name: "Kotobuki",
	}
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: `{"id":"548ab5b00904b8bf2e8dd838","name":"Kotobuki","alias":"kotobuki","users":["alice@example.org"],"owner":"alice@example.org"}`,
	}
	team.client = NewClient(&http.Client{Transport: &transport})
	r := team.remove()
	c.Assert(r, Equals, "Team removed successfully.")
}

func (s *S) TestTeamRemoveWithoutTarget(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current:\n"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	team := &Team{
		Name: "Kotobuki",
	}
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: `{}`,
	}
	team.client = NewClient(&http.Client{Transport: &transport})
	r := team.remove()
	c.Assert(r, Equals, "You have not selected any target as default. For more details, please run `backstage target-set -h`.")
}
