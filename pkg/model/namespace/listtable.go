package namespace

import "github.com/containerum/chkit/pkg/model"

var (
	_ model.TableRenderer = &NamespaceList{}
)

func (_ NamespaceList) TableHeaders() []string {
	return new(Namespace).TableHeaders()
}

func (list NamespaceList) TableRows() [][]string {
	row := make([][]string, 0, len(list))
	for _, ns := range list {
		row = append(row, ns.TableRows()...)
	}
	return row
}
