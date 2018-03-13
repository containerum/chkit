package model

import (
	"bytes"

	"github.com/olekukonko/tablewriter"
)

type TableRenderer interface {
	TableHeaders() []string
	TableRows() [][]string
}

func RenderTable(renderer TableRenderer) string {
	buf := &bytes.Buffer{}
	table := tablewriter.NewWriter(buf)
	table.SetHeader(renderer.TableHeaders())
	table.AppendBulk(renderer.TableRows())
	table.Render()
	return buf.String()
}
