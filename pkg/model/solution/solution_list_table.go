package solution

import "github.com/containerum/chkit/pkg/model"

func (list SolutionsList) RenderTable() string {
	return model.RenderTable(list)
}

func (SolutionsList) TableHeaders() []string {
	return (Solution{}).TableHeaders()
}

func (list SolutionsList) TableRows() [][]string {
	var rows = make([][]string, 0, list.Len())
	for _, solution := range list.Solutions {
		rows = append(rows, SolutionFromKube(solution).TableRows()...)
	}
	return rows
}
