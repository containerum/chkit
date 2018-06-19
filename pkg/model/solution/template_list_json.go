package solution

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = SolutionList{}
)

func (list SolutionList) RenderJSON() (string, error) {
	data, err := list.MarshalJSON()
	return string(data), err
}

func (list SolutionList) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(list, "", model.Indent)
	return data, err
}
