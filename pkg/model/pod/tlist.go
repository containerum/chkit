package pod

import (
	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = new(PodList)
)

func (list PodList) RenderTable() string {
	return model.RenderTable(list)
}

func (_ PodList) TableHeaders() []string {
	return new(Pod).TableHeaders()
}

func (list PodList) TableRows() [][]string {
	table := make([][]string, 0, len(list))
	for _, pod := range list {
		table = append(table, pod.TableRows()...)
	}
	return table
}
