package solution

import (
	"bytes"

	"fmt"

	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/text"
)

var (
	_ model.TableRenderer = UserSolution{}
)

func (solution UserSolution) RenderTable() string {
	return model.RenderTable(solution)
}

func (UserSolution) TableHeaders() []string {
	return []string{
		"Name",
		"Template",
		"Namespace",
		"ENV",
	}
}

func (solution UserSolution) TableRows() [][]string {
	const envWidth = 32
	return [][]string{{
		solution.Name,
		solution.Template,
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
