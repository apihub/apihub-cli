package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

type Table struct {
	Content [][]string
	Header  []string
}

func (t *Table) Render() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(t.Header)

	for _, v := range t.Content {
		table.Append(v)
	}
	table.Render()
}