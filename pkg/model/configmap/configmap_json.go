package configmap

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = ConfigMap{}
	_ json.Marshaler     = ConfigMap{}
)

func (config ConfigMap) RenderJSON() (string, error) {
	data, err := config.MarshalJSON()
	return string(data), err
}

func (config ConfigMap) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(config.ToKube(), "", model.Indent)
}
