package solution

import (
	"bytes"

	"fmt"

	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/text"
)

var (
	_ model.TableRenderer = Solution{}
)

func (solution Solution) RenderTable() string {
	return model.RenderTable(solution)
}

func (Solution) TableHeaders() []string {
	return []string{
		"Name",
		"Template",
		"Branch",
		"Namespace",
		"ENV",
	}
}

func (solution Solution) TableRows() [][]string {
	const envWidth = 48
	return [][]string{{
		solution.Name,
		solution.Template,
		solution.Branch,
		solution.Namespace,
		func() string {
			buf := bytes.NewBuffer(make([]byte, 0, (envWidth+1)*len(solution.Env)))
			for k, v := range solution.Env {
				env := text.Crop(fmt.Sprintf("%s->%q", k, v), envWidth) + "\n"
				buf.WriteString(env)
			}
			return buf.String()
		}(),
	}}
}
