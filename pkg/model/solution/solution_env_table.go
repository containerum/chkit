package solution

import (
	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = Solution{}
)

func (solution SolutionEnv) RenderTable() string {
	return model.RenderTable(solution)
}

func (SolutionEnv) TableHeaders() []string {
	return []string{
		"Key",
		"Value",
	}
}

func (solutionEnv SolutionEnv) TableRows() [][]string {
	var rows = make([][]string, 0, len(solutionEnv.Env))
	var env []string

	for k, v := range solutionEnv.Env {
		env = []string{
			k,
			v,
		}
		rows = append(rows, env)
	}
	return rows
}
