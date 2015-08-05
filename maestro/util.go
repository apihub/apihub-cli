package apihub

import (
	"os"
	"path"
	"sort"
)

func JoinHomePath(p ...string) string {
	ps := []string{os.ExpandEnv("$HOME")}
	ps = append(ps, p...)
	return path.Join(ps...)
}

func SortMapKeys(m map[string]string) []string {
	mk := make([]string, len(m))

	i := 0
	for k := range m {
		mk[i] = k
		i++
	}

	sort.Strings(mk)
	return mk
}
