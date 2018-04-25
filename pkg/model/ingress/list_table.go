package ingress

import "github.com/containerum/chkit/pkg/model"

func (list IngressList) RenderTable() string {
	return model.RenderTable(list)
}

func (IngressList) TableHeaders() []string {
	return Ingress{}.TableHeaders()
}

func (list IngressList) TableRows() [][]string {
	rows := make([][]string, 0, len(list))
	for _, ingr := range list {
		rows = append(rows, ingr.TableRows()...)
	}
	return rows
}
