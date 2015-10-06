package apihub_test

import (
	"github.com/apihub/apihub-cli/apihub"
	"github.com/tsuru/tsuru/fs/fstest"
	. "gopkg.in/check.v1"
)

func (s *S) TestWriteToken(c *C) {
	rfs := &fstest.RecordingFs{}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()
	err := apihub.WriteToken("Token xyz")
	c.Assert(err, IsNil)
	c.Assert(rfs.HasAction("openfile "+apihub.TokenFileName+" with mode 0600"), Equals, true)
}

func (s *S) TestReadToken(c *C) {
	rfs := &fstest.RecordingFs{FileContent: "Token xyz"}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()
	token, err := apihub.ReadToken()
	c.Assert(err, IsNil)
	c.Assert(token, Equals, "Token xyz")
	c.Assert(rfs.HasAction("openfile "+apihub.TokenFileName+" with mode 0600"), Equals, true)
}

func (s *S) TestReadTokenWhenFileNotFound(c *C) {
	rfs := &fstest.FileNotFoundFs{}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()
	_, err := apihub.ReadToken()
	c.Assert(err, Not(IsNil))
}

func (s *S) TestDeleteToken(c *C) {
	rfs := &fstest.RecordingFs{}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()
	err := apihub.DeleteToken()
	c.Assert(err, IsNil)
	c.Assert(rfs.HasAction("remove "+apihub.TokenFileName), Equals, true)
}

func (s *S) TestDeleteTokenWhenFileNotFound(c *C) {
	rfs := &fstest.FileNotFoundFs{}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()
	err := apihub.DeleteToken()
	c.Assert(err, Not(IsNil))
}
