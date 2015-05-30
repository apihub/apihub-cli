package backstage_test

import (
	"github.com/backstage/backstage-client/backstage"
	"github.com/tsuru/tsuru/fs/fstest"
	. "gopkg.in/check.v1"
)

func (s *S) TestWriteToken(c *C) {
	rfs := &fstest.RecordingFs{}
	backstage.Fsystem = rfs
	defer func() {
		backstage.Fsystem = nil
	}()
	err := backstage.WriteToken("Token xyz")
	c.Assert(err, IsNil)
	c.Assert(rfs.HasAction("openfile "+backstage.TokenFileName+" with mode 0600"), Equals, true)
}

func (s *S) TestReadToken(c *C) {
	rfs := &fstest.RecordingFs{FileContent: "Token xyz"}
	backstage.Fsystem = rfs
	defer func() {
		backstage.Fsystem = nil
	}()
	token, err := backstage.ReadToken()
	c.Assert(err, IsNil)
	c.Assert(token, Equals, "Token xyz")
	c.Assert(rfs.HasAction("openfile "+backstage.TokenFileName+" with mode 0600"), Equals, true)
}

func (s *S) TestReadTokenWhenFileNotFound(c *C) {
	rfs := &fstest.FileNotFoundFs{}
	backstage.Fsystem = rfs
	defer func() {
		backstage.Fsystem = nil
	}()
	_, err := backstage.ReadToken()
	c.Assert(err, Not(IsNil))
}

func (s *S) TestDeleteToken(c *C) {
	rfs := &fstest.RecordingFs{}
	backstage.Fsystem = rfs
	defer func() {
		backstage.Fsystem = nil
	}()
	err := backstage.DeleteToken()
	c.Assert(err, IsNil)
	c.Assert(rfs.HasAction("remove "+backstage.TokenFileName), Equals, true)
}

func (s *S) TestDeleteTokenWhenFileNotFound(c *C) {
	rfs := &fstest.FileNotFoundFs{}
	backstage.Fsystem = rfs
	defer func() {
		backstage.Fsystem = nil
	}()
	err := backstage.DeleteToken()
	c.Assert(err, Not(IsNil))
}
