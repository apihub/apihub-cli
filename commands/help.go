package commands

import (
	"fmt"
	"regexp"
)

func Confirm(ctx *Context, question string) bool {
	r, _ := regexp.Compile("[Y|y]")
	fmt.Fprintf(ctx.Stdout, "%s (y/n)", question)
	var answer string
	fmt.Fscanf(ctx.Stdin, "%s", &answer)
	if !r.MatchString(answer) {
		return false
	}
	return true
}
