package main

import "github.com/olekukonko/tablewriter"

type Table struct {
	Content [][]string
	Header  []string
}

func (t *Table) Render(context *Context) {
	table := tablewriter.NewWriter(context.Stdout)
	table.SetHeader(t.Header)

	for _, v := range t.Content {
		table.Append(v)
	}
	table.Render()
}