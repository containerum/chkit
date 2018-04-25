package deployment

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = Deployment{}
	_ json.Marshaler     = Deployment{}
)

func (depl Deployment) RenderJSON() (string, error) {
	depl.ToKube()
	data, err := depl.MarshalJSON()
	return string(data), err
}

func (depl Deployment) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(depl.origin, "", model.Indent)
	return data, err
}
