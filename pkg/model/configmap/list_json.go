package configmap

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = ConfigMapList{}
	_ json.Marshaler     = ConfigMap{}
)

func (list ConfigMapList) RenderJSON() (string, error) {
	data, err := list.MarshalJSON()
	return string(data), err
}

func (list ConfigMapList) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(list.ToKube(), "", model.Indent)
}
