package apihub_test

import (
	"os"

	"github.com/apihub/apihub-cli/apihub"
	. "gopkg.in/check.v1"
)

func (s *S) TestJoinHomePath(c *C) {
	str := ".apihub_targets"
	home := os.ExpandEnv("$HOME")
	c.Assert(apihub.JoinHomePath(str), Equals, home+"/"+str)
}

func (s *S) TestJoinHomePathWithMultipleValues(c *C) {
	str := ".apihub_targets"
	str2 := ".test"
	home := os.ExpandEnv("$HOME")
	c.Assert(apihub.JoinHomePath(str, str2), Equals, home+"/"+str+"/"+str2)
}

func (s *S) TestSortMapKeys(c *C) {
	mapkeys := map[string]string{"c": "c", "b": "b", "a": "a"}
	sortedKeys := apihub.SortMapKeys(mapkeys)
	c.Assert(sortedKeys, DeepEquals, []string{"a", "b", "c"})
}
