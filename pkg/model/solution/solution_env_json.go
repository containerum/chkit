package solution

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = SolutionEnv{}
)

func (envs SolutionEnv) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(envs, "", model.Indent)
	return string(data), err
}
