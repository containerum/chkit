package solution

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = SolutionsList{}
)

func (list SolutionsList) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(list, "", model.Indent)
	return string(data), err
}
