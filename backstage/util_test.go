package main

import (
	"os"

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
