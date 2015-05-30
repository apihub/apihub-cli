package commands

import (
	"fmt"

	"github.com/olekukonko/tablewriter"
)

type Table struct {
	Title   string
	Content [][]string
	Header  []string
}

func (t *Table) Render(context *Context) {
	if t.Title != "" {
		fmt.Println(t.Title)
	}
	table := tablewriter.NewWriter(context.Stdout)
	table.SetHeader(t.Header)

	for _, v := range t.Content {
		table.Append(v)
	}
	table.Render()
}
