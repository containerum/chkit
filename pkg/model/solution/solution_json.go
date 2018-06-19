package solution

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = UserSolution{}
	_ json.Marshaler     = UserSolution{}
)

func (solution UserSolution) RenderJSON() (string, error) {
	data, err := solution.MarshalJSON()
	return string(data), err
}

func (solution UserSolution) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(solution.ToKube(), "", model.Indent)
}
