package configmap

import (
	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = ConfigMapList{}
)

func (list ConfigMapList) RenderTable() string {
	return model.RenderTable(list)
}

func (ConfigMapList) TableHeaders() []string {
	return ConfigMap{}.TableHeaders()
}

func (list ConfigMapList) TableRows() [][]string {
	var rows = make([][]string, 0, list.Len())
	for _, config := range list {
		rows = append(rows, config.TableRows()...)
	}
	return rows
}
