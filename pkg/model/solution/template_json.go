package solution

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = Solution{}
	_ json.Marshaler     = Solution{}
)

func (solution Solution) RenderJSON() (string, error) {
	data, err := solution.MarshalJSON()
	return string(data), err
}

func (solution Solution) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(solution.ToKube(), "", model.Indent)
	return data, err
}
