package main

import (
	"os"

	"github.com/tsuru/tsuru/fs/testing"
	. "gopkg.in/check.v1"
)

func (s *S) TestJoinHomePath(c *C) {
	str := ".backstage_targets"
	home := os.ExpandEnv("$HOME")
	c.Assert(joinHomePath(str), Equals, home+"/"+str)
}

func (s *S) TestJoinHomePathWithMultipleValues(c *C) {
	str := ".backstage_targets"
	str2 := ".test"
	home := os.ExpandEnv("$HOME")
	c.Assert(joinHomePath(str, str2), Equals, home+"/"+str+"/"+str2)
}

func (s *S) TestLoginRequired(c *C) {
	rfs := &testing.RecordingFs{FileContent: "Token xyz"}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	login := LoginRequired()
	c.Assert(login, IsNil)
}

func (s *S) TestLoginRequriedWhenFileNotFound(c *C) {
	rfs := &testing.FileNotFoundFs{}
	fsystem = rfs
	defer func() {
		fsystem = nil
	}()
	login := LoginRequired()
	c.Assert(login, Not(IsNil))
}
