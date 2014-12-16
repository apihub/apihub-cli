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
		Alias: "kotobuki",
	}
	transport := ttesting.Transport{
		Status:  http.StatusCreated,
		Message: `{"id":"548ab5b00904b8bf2e8dd838","name":"Kotobuki","alias":"kotobuki","users":["alice@example.org"],"owner":"alice@example.org"}`,
	}
	team.client = NewClient(&http.Client{Transport: &transport})
	r := team.save()
	c.Assert(r, Equals, "Your team has been created.")
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
		Alias: "kotobuki",
	}
	team.client = NewClient(&http.Client{Transport: &transport})
	r := team.save()
	c.Assert(r, Equals, "Someone already has that team name. Could you try another?")
}

func (s *S) TestTeamList(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: `[{"id":"54825cd18f897dbba8aba570","name":"backstage","users":["alice@example.org"],"owner":"alice@example.org"}]`,
	}
	team := &Team{}
	team.client = NewClient(&http.Client{Transport: &transport})
	table, err := team.list()
	c.Assert(err, IsNil)
	c.Assert(table.Header, DeepEquals, []string{"Team Name", "Alias", "Owner"})
	c.Assert(table.Content, DeepEquals, [][]string{[]string{"backstage", "", "alice@example.org"}})
}

func (s *S) TestTeamListWithoutTeam(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: `[]`,
	}
	team := &Team{}
	team.client = NewClient(&http.Client{Transport: &transport})
	table, err := team.list()
	c.Assert(table, IsNil)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "You have no teams.")
}

func (s *S) TestTeamInfo(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: `{"id":"54825cd18f897dbba8aba570","name":"Backstage","alias":"backstage","users":["alice@example.org"],"owner":"alice@example.org"}`,
	}
	team := &Team{}
	team.client = NewClient(&http.Client{Transport: &transport})
	table, err := team.info()
	c.Assert(err, IsNil)
	c.Assert(table.Header, DeepEquals, []string{"Team Members"})
	c.Assert(table.Content, DeepEquals, [][]string{[]string{"alice@example.org"}})
}

func (s *S) TestTeamRemove(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	team := &Team{
		Alias: "kotobuki",
	}
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: `{"id":"548ab5b00904b8bf2e8dd838","name":"Kotobuki","alias":"kotobuki","users":["alice@example.org"],"owner":"alice@example.org"}`,
	}
	team.client = NewClient(&http.Client{Transport: &transport})
	r := team.remove()
	c.Assert(r, Equals, "Your team has been deleted.")
}

func (s *S) TestTeamRemoveWithoutTarget(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current:\n"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	team := &Team{
		Alias: "kotobuki",
	}
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: `{}`,
	}
	team.client = NewClient(&http.Client{Transport: &transport})
	r := team.remove()
	c.Assert(r, Equals, "You have not selected any target as default. For more details, please run `backstage target-set -h`.")
}

func (s *S) TestTeamAddUser(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	team := &Team{
		Alias: "kotobuki",
	}
	transport := ttesting.Transport{
		Status:  http.StatusCreated,
		Message: `{"id":"548ab5b00904b8bf2e8dd838","name":"Kotobuki","alias":"kotobuki","users":["alice@example.org"],"owner":"alice@example.org"}`,
	}
	team.client = NewClient(&http.Client{Transport: &transport})
	r := team.addUser("alice@example.org")
	c.Assert(r, Equals, "User `alice@example.org` added successfully to team `kotobuki`.")
}

func (s *S) TestTeamAddUserWhenUserDoesNotExist(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	team := &Team{
		Alias: "kotobuki",
	}
	transport := ttesting.Transport{
		Status:  http.StatusCreated,
		Message: `{"id":"548ab5b00904b8bf2e8dd838","name":"Kotobuki","alias":"kotobuki","users":["alice@example.org"],"owner":"alice@example.org"}`,
	}
	team.client = NewClient(&http.Client{Transport: &transport})
	r := team.addUser("invalid-email@example.org")
	c.Assert(r, Equals, "Sorry, the user was not found.")
}

func (s *S) TestTeamRemoveUser(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	team := &Team{
		Alias: "kotobuki",
	}
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: `{"id":"548ab5b00904b8bf2e8dd838","name":"Kotobuki","alias":"kotobuki","users":["alice@example.org"],"owner":"alice@example.org"}`,
	}
	team.client = NewClient(&http.Client{Transport: &transport})
	r := team.removeUser("ringo@example.org")
	c.Assert(r, Equals, "User `ringo@example.org` removed successfully to team `kotobuki`.")
}

func (s *S) TestTeamRemoveUserWhenUserItTheOwner(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	team := &Team{
		Alias: "kotobuki",
	}
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: `{"id":"548ab5b00904b8bf2e8dd838","name":"Kotobuki","alias":"kotobuki","users":["alice@example.org"],"owner":"alice@example.org"}`,
	}
	team.client = NewClient(&http.Client{Transport: &transport})
	r := team.removeUser("alice@example.org")
	c.Assert(r, Equals, "It's not allowed to remove the owner from its team.")
}
