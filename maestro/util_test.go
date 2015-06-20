package backstage_test

import (
	"os"

	"github.com/backstage/backstage-cli/maestro"
	. "gopkg.in/check.v1"
)

func (s *S) TestJoinHomePath(c *C) {
	str := ".backstage_targets"
	home := os.ExpandEnv("$HOME")
	c.Assert(backstage.JoinHomePath(str), Equals, home+"/"+str)
}

func (s *S) TestJoinHomePathWithMultipleValues(c *C) {
	str := ".backstage_targets"
	str2 := ".test"
	home := os.ExpandEnv("$HOME")
	c.Assert(backstage.JoinHomePath(str, str2), Equals, home+"/"+str+"/"+str2)
}

func (s *S) TestSortMapKeys(c *C) {
	mapkeys := map[string]string{"c": "c", "b": "b", "a": "a"}
	sortedKeys := backstage.SortMapKeys(mapkeys)
	c.Assert(sortedKeys, DeepEquals, []string{"a", "b", "c"})
}
