package solution

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = SolutionTemplate{}
	_ json.Marshaler     = SolutionTemplate{}
)

func (solution SolutionTemplate) RenderJSON() (string, error) {
	data, err := solution.MarshalJSON()
	return string(data), err
}

func (solution SolutionTemplate) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(solution.ToKube(), "", model.Indent)
	return data, err
}
