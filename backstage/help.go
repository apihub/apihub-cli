package main

import (
	"fmt"
	"regexp"
)

var AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} command{{if .Flags}} [command options]{{end}} [arguments...]

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`

func Confirm(ctx *Context, question string) bool {
	r, _ := regexp.Compile("[Y|y]")
	fmt.Fprintf(ctx.Stdout, "%s (y/n)", question)
	var answer string
	fmt.Fscanf(ctx.Stdin, "%s", &answer)
	if !r.MatchString(answer) {
		fmt.Fprintln(ctx.Stdout, "Operation cancelled.")
		return false
	}
	return true
}
