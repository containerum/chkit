package model

import (
	"bytes"

	"github.com/olekukonko/tablewriter"
)

type TableItem interface {
	TableHeaders() []string
	TableRows() [][]string
}
type TableRenderer interface {
	TableItem
	RenderTable() string
}

func RenderTable(renderer TableItem) string {
	buf := &bytes.Buffer{}
	table := tablewriter.NewWriter(buf)
	table.SetHeader(renderer.TableHeaders())
	table.AppendBulk(renderer.TableRows())
	table.Render()
	return buf.String()
}
