package namespace

import (
	"strconv"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = &NamespaceList{}
)

func (_ NamespaceList) TableHeaders() []string {
	return append([]string{"№"}, new(Namespace).TableHeaders()...)
}

func (list NamespaceList) TableRows() [][]string {
	rows := make([][]string, 0, len(list))
	for i, ns := range list {
		var nsRows = ns.TableRows()
		for _, nsRows := range nsRows {
			rows = append(rows, append([]string{strconv.Itoa(i + 1)}, nsRows...))
		}
	}
	return rows
}

func (list NamespaceList) RenderTable() string {
	return model.RenderTable(list)
}
