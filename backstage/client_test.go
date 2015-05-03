package main

import (
	"net/http"

	"github.com/tsuru/tsuru/cmd/cmdtest"
	"github.com/tsuru/tsuru/fs/fstest"
	. "gopkg.in/check.v1"
)

func (s *S) TestClientCreate(c *C) {
	rfs := &fstest.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	client := &Client{
		Id:          "backstage",
		Name:        "Backstage",
		RedirectUri: "http://www.example.org/auth",
	}
	transport := cmdtest.Transport{
		Status:  http.StatusCreated,
		Message: `{"id":"backstage","secret":"TJl5HvdhC-NepxCAUXy7fanL4enr3xKDiUcWI2KrBSY=","name":"Backstage","redirect_uri":"","owner":"owner@example.org","team":"backstage"}`,
	}
	client.client = NewHTTPClient(&http.Client{Transport: &transport})
	r := client.save()
	c.Assert(r, Equals, "Your new client has been created.")
}

func (s *S) TestClientCreateWithInvalidTeam(c *C) {
	rfs := &fstest.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	transport := cmdtest.Transport{
		Status:  http.StatusBadRequest,
		Message: `{"error":"bad_request","error_description":"Team not found."}`,
	}
	client := &Client{
		Team: "backstage",
	}
	client.client = NewHTTPClient(&http.Client{Transport: &transport})
	r := client.save()
	c.Assert(r, Equals, "Team not found.")
}

func (s *S) TestClientCreateWithAnExistingName(c *C) {
	rfs := &fstest.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	transport := cmdtest.Transport{
		Status:  http.StatusBadRequest,
		Message: `{"error":"bad_request","error_description":"There is another client with this name."}`,
	}
	client := &Client{
		Name: "backstage",
	}
	client.client = NewHTTPClient(&http.Client{Transport: &transport})
	r := client.save()
	c.Assert(r, Equals, "There is another client with this name.")
}

func (s *S) TestClientRemove(c *C) {
	rfs := &fstest.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	client := &Client{
		Id: "backstage",
	}
	transport := cmdtest.Transport{
		Status:  http.StatusOK,
		Message: `{"id":"backstage","secret":"TJl5HvdhC-NepxCAUXy7fanL4enr3xKDiUcWI2KrBSY=","name":"Backstage","redirect_uri":"","owner":"owner@example.org","team":"backstage"}`,
	}
	client.client = NewHTTPClient(&http.Client{Transport: &transport})
	r := client.remove()
	c.Assert(r, Equals, "The client `backstage` has been deleted.")
}

func (s *S) TestClientRemoveWithError(c *C) {
	rfs := &fstest.RecordingFs{FileContent: "current:\n"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	client := &Client{
		Name: "backstage",
	}
	transport := cmdtest.Transport{
		Status:  http.StatusOK,
		Message: `{}`,
	}
	client.client = NewHTTPClient(&http.Client{Transport: &transport})
	r := client.remove()
	c.Assert(r, Equals, "You have not selected any target as default. For more details, please run `backstage target-set -h`.")
}
