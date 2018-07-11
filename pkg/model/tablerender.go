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
	RenderTable() string
}

func RenderTable(renderer TableItem) string {
	buf := &bytes.Buffer{}

	table := tablewriter.NewWriter(buf)
	table.SetAutoWrapText(false)
	table.SetRowSeparator("_")
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetReflowDuringAutoWrap(true)
	table.SetCenterSeparator("_")
	table.SetColumnSeparator(" ")
	table.SetHeader(renderer.TableHeaders())
	table.AppendBulk(renderer.TableRows())
	table.Render()
	return buf.String()
}
