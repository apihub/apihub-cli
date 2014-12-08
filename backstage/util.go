package main

import (
	"os"
	"path"
)

func joinHomePath(p ...string) string {
	ps := []string{os.ExpandEnv("$HOME")}
	ps = append(ps, p...)
	return path.Join(ps...)
}
