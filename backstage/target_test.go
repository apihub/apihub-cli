package main

import (
	"os"
	"path"

	"github.com/tsuru/tsuru/fs/testing"
	. "gopkg.in/check.v1"
)

func (s *S) TestLoadTargets(c *C) {
	rfs := &testing.RecordingFs{}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	_, err := LoadTargets()
	c.Assert(err, IsNil)
	filePath := path.Join(os.ExpandEnv("${HOME}"), ".backstage_targets")
	c.Assert(rfs.HasAction("openfile "+filePath+" with mode 0600"), Equals, true)
}

func (s *S) TestAddNewTarget(c *C) {
	rfs := &testing.RecordingFs{}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	t, err := LoadTargets()
	c.Assert(err, IsNil)
	err = t.add("backstage", "http://www.example.org")
	c.Assert(err, IsNil)
	t, _ = LoadTargets()
	c.Assert(t.Options["backstage"], Equals, "http://www.example.org")
}

func (s *S) TestAddWhenLabelAlreadyExists(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	t, err := LoadTargets()
	err = t.add("backstage", "http://www.example.org")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Sorry, that label has been used by another user.")

}
func (s *S) TestListTargets(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com\n  example: www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	t, err := LoadTargets()
	c.Assert(err, IsNil)
	targetTable := t.list()
	c.Assert(targetTable.Header, DeepEquals, []string{"Default", "Label", "Backstage Server"})
	c.Assert(targetTable.Content, DeepEquals, [][]string{[]string{"*", "backstage", "http://www.example.com"}, []string{"", "example", "www.example.com"}})
}

func (s *S) TestRemoveTarget(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	t, err := LoadTargets()
	c.Assert(err, IsNil)
	err = t.remove("backstage")
	c.Assert(err, IsNil)
}

func (s *S) TestRemoveTargetWithInvalidLabel(c *C) {
	rfs := &testing.RecordingFs{}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	t, err := LoadTargets()
	c.Assert(err, IsNil)
	err = t.remove("invalid-label")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Sorry, that label does not exist.")
}

func (s *S) TestSetTargetAsDefault(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current:\noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	t, err := LoadTargets()
	c.Assert(t.Current, Equals, "")
	c.Assert(err, IsNil)
	err = t.setDefault("backstage")
	c.Assert(err, IsNil)
	t, err = LoadTargets()
	c.Assert(t.Current, Equals, "backstage")
}

func (s *S) TestSetTargetAsDefaultWithInvalidLabel(c *C) {
	rfs := &testing.RecordingFs{}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	t, err := LoadTargets()
	c.Assert(err, IsNil)
	err = t.setDefault("invalid-label")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Sorry, that label does not exist.")
}

func (s *S) TestGetURL(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  backstage: http://www.example.com/"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	url, err := GetURL("/api/teams")
	c.Assert(err, IsNil)
	c.Assert(url, Equals, "http://www.example.com/api/teams")
}

func (s *S) TestGetURLWithoutEndpoint(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: backstage\noptions:\n  key: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	url, err := GetURL("/api/teams")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "You have not selected any target as default. For more details, please run `backstage target-set -h`.")
	c.Assert(url, Equals, "")
}

func (s *S) TestGetURLWithoutCurrent(c *C) {
	rfs := &testing.RecordingFs{FileContent: "current: \noptions:\n  backstage: http://www.example.com"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	url, err := GetURL("/api/teams")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "You have not selected any target as default. For more details, please run `backstage target-set -h`.")
	c.Assert(url, Equals, "")
}

func (s *S) TestGetURLWithoutContent(c *C) {
	rfs := &testing.RecordingFs{}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	url, err := GetURL("/api/teams")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "You have not selected any target as default. For more details, please run `backstage target-set -h`.")
	c.Assert(url, Equals, "")
}