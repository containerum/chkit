package solution

import (
	"strings"

	"github.com/containerum/chkit/pkg/model"
)

func (solution Solution) RenderTable() string {
	return model.RenderTable(solution)
}

func (Solution) TableHeaders() []string {
	return []string{
		"Name",
		"URL",
		"Images",
	}
}

func (solution Solution) TableRows() [][]string {
	return [][]string{{
		solution.Name,
		solution.URL,
		strings.Join(solution.Images, "\n"),
	}}
}
