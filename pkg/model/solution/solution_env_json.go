package solution

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = SolutionEnv{}
	_ json.Marshaler     = SolutionEnv{}
)

func (solution SolutionEnv) RenderJSON() (string, error) {
	data, err := solution.MarshalJSON()
	return string(data), err
}

func (solution SolutionEnv) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(solution, "", model.Indent)
}
