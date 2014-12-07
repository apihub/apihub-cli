package main

import (
	"bytes"
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) TestConfirmCommandReturnTrueFory(c *C) {
	var stdout bytes.Buffer
	context := Context{Stdout: &stdout, Stdin: strings.NewReader("y\n")}
	result := Confirm(&context, "Are you sure you want to delete it?")
	c.Assert(result, Equals, true)
	c.Assert(stdout.String(), Equals, "Are you sure you want to delete it? (y/n)")
}

func (s *S) TestConfirmCommandReturnTrueForY(c *C) {
	var stdout bytes.Buffer
	context := Context{Stdout: &stdout, Stdin: strings.NewReader("Y\n")}
	result := Confirm(&context, "Are you sure you want to delete it?")
	c.Assert(result, Equals, true)
	c.Assert(stdout.String(), Equals, "Are you sure you want to delete it? (y/n)")
}

func (s *S) TestConfirmCommandReturnFalseForn(c *C) {
	var stdout bytes.Buffer
	context := Context{Stdout: &stdout, Stdin: strings.NewReader("n\n")}
	result := Confirm(&context, "Are you sure you want to delete it?")
	c.Assert(result, Equals, false)
	c.Assert(stdout.String(), Equals, "Are you sure you want to delete it? (y/n)Operation cancelled.\n")
}

func (s *S) TestConfirmCommandReturnFalseForN(c *C) {
	var stdout bytes.Buffer
	context := Context{Stdout: &stdout, Stdin: strings.NewReader("N\n")}
	result := Confirm(&context, "Are you sure you want to delete it?")
	c.Assert(result, Equals, false)
	c.Assert(stdout.String(), Equals, "Are you sure you want to delete it? (y/n)Operation cancelled.\n")
}

func (s *S) TestConfirmCommandReturnFalseForAnythingButY(c *C) {
	var stdout bytes.Buffer
	context := Context{Stdout: &stdout, Stdin: strings.NewReader("A\n")}
	result := Confirm(&context, "Are you sure you want to delete it?")
	c.Assert(result, Equals, false)
	c.Assert(stdout.String(), Equals, "Are you sure you want to delete it? (y/n)Operation cancelled.\n")
}
