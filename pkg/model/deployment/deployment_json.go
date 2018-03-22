package deployment

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = Deployment{}
)

func (depl Deployment) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(depl, "", "    ")
	return string(data), err
}
