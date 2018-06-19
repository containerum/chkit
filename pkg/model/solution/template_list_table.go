package solution

import "github.com/containerum/chkit/pkg/model"

func (solution TemplatesList) RenderTable() string {
	return model.RenderTable(solution)
}

func (TemplatesList) TableHeaders() []string {
	return (SolutionTemplate{}).TableHeaders()
}

func (list TemplatesList) TableRows() [][]string {
	var rows = make([][]string, 0, list.Len())
	for _, solution := range list.Solutions {
		rows = append(rows, SolutionTemplateFromKube(solution).TableRows()...)
	}
	return rows
}
