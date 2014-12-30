package main

import (
	"net/http"

	ttesting "github.com/tsuru/tsuru/cmd/testing"
	"github.com/tsuru/tsuru/fs/testing"
	. "gopkg.in/check.v1"
)

func (s *S) TestServiceCreate(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	service := &Service{
		Subdomain:       "backstage",
		AllowKeylessUse: true,
		Description:     "test",
		Disabled:        false,
		Documentation:   "http://www.example.org/doc",
		Endpoint:        "http://github.com/backstage",
		Owner:           "alice@example.org",
		Timeout:         10,
	}
	transport := ttesting.Transport{
		Status:  http.StatusCreated,
		Message: `{"subdomain":"backstage","created_at":"2014-12-05T17:44:39.462-02:00","updated_at":"2014-12-05T17:44:39.462-02:00","allow_keyless_use":true,"description":"test","disabled":false,"documentation":"http://www.example.org/doc","endpoint":"http://github.com/backstage","owner":"alice@example.org","timeout":10}`,
	}
	service.client = NewClient(&http.Client{Transport: &transport})
	r := service.save()
	c.Assert(r, Equals, "Your new service has been created.")
}

func (s *S) TestServiceCreateWithInvalidSubdomain(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	transport := ttesting.Transport{
		Status:  http.StatusBadRequest,
		Message: `{"status_code":400,"error_description":"Service not found."}`,
	}
	service := &Service{
		Subdomain: "backstage",
	}
	service.client = NewClient(&http.Client{Transport: &transport})
	r := service.save()
	c.Assert(r, Equals, "Service not found.")
}

func (s *S) TestServiceCreateWithAnExistingSubdomain(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	transport := ttesting.Transport{
		Status:  http.StatusBadRequest,
		Message: `{"status_code":400,"error_description":"There is another service with this subdomain."}`,
	}
	service := &Service{
		Subdomain: "backstage",
	}
	service.client = NewClient(&http.Client{Transport: &transport})
	r := service.save()
	c.Assert(r, Equals, "There is another service with this subdomain.")
}

func (s *S) TestServiceRemove(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	service := &Service{
		Subdomain: "backstage",
	}
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: `{"subdomain":"backstage","created_at":"2014-12-05T17:44:39.462-02:00","updated_at":"2014-12-05T17:44:39.462-02:00","allow_keyless_use":true,"description":"test","disabled":false,"documentation":"http://www.example.org/doc","endpoint":"http://github.com/backstage","owner":"alice@example.org","timeout":10}`,
	}
	service.client = NewClient(&http.Client{Transport: &transport})
	r := service.remove()
	c.Assert(r, Equals, "The service `backstage` has been deleted.")
}

func (s *S) TestServiceRemoveWithInvalidSubdomain(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current:\n"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	service := &Service{
		Subdomain: "backstage",
	}
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: `{}`,
	}
	service.client = NewClient(&http.Client{Transport: &transport})
	r := service.remove()
	c.Assert(r, Equals, "You have not selected any target as default. For more details, please run `backstage target-set -h`.")
}
