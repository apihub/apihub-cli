package apihub_test

import (
	"os"
	"path"
	"strings"

	"github.com/apihub/apihub-cli/maestro"
	"github.com/tsuru/tsuru/fs/fstest"
	. "gopkg.in/check.v1"
)

func (s *S) TestLoadTargets(c *C) {
	rfs := &fstest.RecordingFs{}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()
	_, err := apihub.LoadTargets()
	c.Assert(err, IsNil)
	filePath := path.Join(os.ExpandEnv("${HOME}"), ".apihub_targets")
	c.Assert(rfs.HasAction("openfile "+filePath+" with mode 0600"), Equals, true)
}

func (s *S) TestAddNewTarget(c *C) {
	rfs := &fstest.RecordingFs{}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()
	t, err := apihub.LoadTargets()
	c.Assert(err, IsNil)
	err = t.Add("apihub", "http://www.example.org")
	c.Assert(err, IsNil)
	t, _ = apihub.LoadTargets()
	c.Assert(t.Options["apihub"], Equals, "http://www.example.org")
}

func (s *S) TestAddWhenLabelAlreadyExists(c *C) {
	rfs := &fstest.RecordingFs{FileContent: "current: apihub\noptions:\n  apihub: http://www.example.com"}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()
	t, err := apihub.LoadTargets()
	err = t.Add("apihub", "http://www.example.org")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Sorry, that label has been used by another user.")

}
func (s *S) TestListTargets(c *C) {
	rfs := &fstest.RecordingFs{FileContent: "current: apihub\noptions:\n  apihub: http://www.example.com\n  example: www.example.com"}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()
	t, err := apihub.LoadTargets()
	c.Assert(err, IsNil)

	current, sortedMapKeys, options := t.GetOptions()
	c.Assert(current, Equals, "apihub")
	c.Assert(sortedMapKeys, DeepEquals, []string{"apihub", "example"})
	c.Assert(options, DeepEquals, map[string]string{
		"apihub":  "http://www.example.com",
		"example": "www.example.com",
	})
}

func (s *S) TestRemoveTarget(c *C) {
	rfs := &fstest.RecordingFs{FileContent: "current: apihub\noptions:\n  apihub: http://www.example.com"}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()
	t, err := apihub.LoadTargets()
	c.Assert(err, IsNil)
	err = t.Remove("apihub")
	c.Assert(err, IsNil)
}

func (s *S) TestRemoveTargetWithInvalidLabel(c *C) {
	rfs := &fstest.RecordingFs{}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()
	t, err := apihub.LoadTargets()
	c.Assert(err, IsNil)
	err = t.Remove("invalid-label")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Target not found.")
}

func (s *S) TestSetTargetAsDefault(c *C) {
	rfs := &fstest.RecordingFs{FileContent: "current:\noptions:\n  apihub: http://www.example.com"}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()
	t, err := apihub.LoadTargets()
	c.Assert(t.Current, Equals, "")
	c.Assert(err, IsNil)
	err = t.SetDefault("apihub")
	c.Assert(err, IsNil)
	t, err = apihub.LoadTargets()
	c.Assert(t.Current, Equals, "apihub")
}

func (s *S) TestSetTargetAsDefaultWithInvalidLabel(c *C) {
	rfs := &fstest.RecordingFs{}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()
	t, err := apihub.LoadTargets()
	c.Assert(err, IsNil)
	err = t.SetDefault("invalid-label")
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Target not found.")
}

func (s *S) TestGetCurrentTarget(c *C) {
	rfs := &fstest.RecordingFs{FileContent: "current: apihub\noptions:\n  apihub: http://www.example.com/"}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()

	currentTarget, err := apihub.GetCurrentTarget()
	url := strings.TrimRight(currentTarget, "/") + "/api/teams"
	c.Assert(err, IsNil)
	c.Assert(url, Equals, "http://www.example.com/api/teams")
}

func (s *S) TestGetCurrentTargetWithoutEndpoint(c *C) {
	rfs := &fstest.RecordingFs{FileContent: "current: apihub\noptions:\n  key: http://www.example.com"}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()

	currentTarget, err := apihub.GetCurrentTarget()
	url := strings.TrimRight(currentTarget, "/") + "/api/teams"
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "You have not selected any target as default. For more details, please run `apihub target-set -h`.")
	c.Assert(url, Equals, "/api/teams")
}

func (s *S) TestGetCurrentTargetWithoutCurrent(c *C) {
	rfs := &fstest.RecordingFs{FileContent: "current: \noptions:\n  apihub: http://www.example.com"}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()

	currentTarget, err := apihub.GetCurrentTarget()
	url := strings.TrimRight(currentTarget, "/") + "/api/teams"
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "You have not selected any target as default. For more details, please run `apihub target-set -h`.")
	c.Assert(url, Equals, "/api/teams")
}

func (s *S) TestGetCurrentTargetWithoutContent(c *C) {
	rfs := &fstest.RecordingFs{}
	apihub.Fsystem = rfs
	defer func() {
		apihub.Fsystem = nil
	}()

	currentTarget, err := apihub.GetCurrentTarget()
	url := strings.TrimRight(currentTarget, "/") + "/api/teams"
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "You have not selected any target as default. For more details, please run `apihub target-set -h`.")
	c.Assert(url, Equals, "/api/teams")
}
