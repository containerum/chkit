package solution

import "github.com/containerum/chkit/pkg/model"

func (solution SolutionList) RenderTable() string {
	return model.RenderTable(solution)
}

func (SolutionList) TableHeaders() []string {
	return (Solution{}).TableHeaders()
}

func (list SolutionList) TableRows() [][]string {
	var rows = make([][]string, 0, list.Len())
	for _, solution := range list.Solutions {
		rows = append(rows, SolutionFromKube(solution).TableRows()...)
	}
	return rows
}
