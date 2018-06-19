package solution

import (
	"strings"

	"github.com/containerum/chkit/pkg/model"
)

func (solution SolutionTemplate) RenderTable() string {
	return model.RenderTable(solution)
}

func (SolutionTemplate) TableHeaders() []string {
	return []string{
		"Name",
		"URL",
		"Images",
	}
}

func (solution SolutionTemplate) TableRows() [][]string {
	return [][]string{{
		solution.Name,
		solution.URL,
		strings.Join(solution.Images, "\n"),
	}}
}
