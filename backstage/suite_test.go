package main

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) TestSetUp(c *C) {
	TargetFileName = "/tmp/.backstage_targets"
}
