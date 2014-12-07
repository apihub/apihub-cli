package main

import (
	"io"

	"github.com/codegangsta/cli"
)

type Context struct {
	cli.Context
	Stdin  io.Reader
	Stdout io.Writer
}
