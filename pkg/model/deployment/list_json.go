package deployment

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = DeploymentList{}
)

func (list DeploymentList) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(list, "", model.Indent)
	return string(data), err
}
