package main

import (
	"net/http"

	ttesting "github.com/tsuru/tsuru/cmd/testing"
	"github.com/tsuru/tsuru/fs/testing"
	. "gopkg.in/check.v1"
)

func (s *S) TestShouldSetCloseToTrue(c *C) {
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, IsNil)
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: "OK",
	}
	client := NewHTTPClient(&http.Client{Transport: &transport})
	client.Do(request)
	c.Assert(request.Close, Equals, true)
}

func (s *S) TestShouldReturnErrorWhenServerIsDown(c *C) {
	rfs := &testing.RecordingFs{FileContent: "http://www.example.org"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, IsNil)
	client := NewHTTPClient(&http.Client{})
	_, err = client.Do(request)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Failed to connect to Backstage server: unsupported protocol scheme \"\"")
}

func (s *S) TestShouldNotIncludeTheHeaderAuthorizationWhenTokenFileIsMissing(c *C) {
	fsystem = &testing.FileNotFoundFs{}
	defer func() {
		fsystem = nil
	}()
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, IsNil)
	trans := ttesting.Transport{
		Message: "",
		Status: http.StatusOK,
	}
	client := NewHTTPClient(&http.Client{Transport: &trans})
	_, err = client.Do(request)
	c.Assert(err, IsNil)
	header := map[string][]string(request.Header)
	_, ok := header["Authorization"]
	c.Assert(ok, Equals, false)
}

func (s *S) TestShouldIncludeTheHeaderAuthorizationWhenTokenFileExists(c *C) {
	fsystem = &testing.RecordingFs{FileContent: "Token mytoken"}
	defer func() {
		fsystem = nil
	}()
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, IsNil)
	trans := ttesting.Transport{
		Message: "",
		Status: http.StatusOK,
	}
	client := NewHTTPClient(&http.Client{Transport: &trans})
	_, err = client.Do(request)
	c.Assert(err, IsNil)
	c.Assert(request.Header.Get("Authorization"), Equals, "Token mytoken")
}

func (s *S) TestShouldIncludeTheClientVersionInTheHeader(c *C) {
	fsystem = &testing.RecordingFs{FileContent: "Token mytoken"}
	defer func() {
		fsystem = nil
	}()
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, IsNil)
	trans := ttesting.Transport{
		Message: "",
		Status: http.StatusOK,
	}
	client := NewHTTPClient(&http.Client{Transport: &trans})
	_, err = client.Do(request)
	c.Assert(err, IsNil)
	c.Assert(request.Header.Get("BackstageClient-Version"), Equals, BackstageClientVersion)
}