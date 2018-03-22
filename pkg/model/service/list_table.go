package service

import "github.com/containerum/chkit/pkg/model"

var (
	_ model.TableRenderer = ServiceList{}
)

func (list ServiceList) RenderTable() string {
	return model.RenderTable(list)
}

func (_ ServiceList) TableHeaders() []string {
	return new(Service).TableHeaders()
}

func (list ServiceList) TableRows() [][]string {
	table := make([][]string, 0, len(list))
	for _, serv := range list {
		table = append(table, serv.TableRows()...)
	}
	return table
}
