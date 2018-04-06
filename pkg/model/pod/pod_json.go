package pod

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = Pod{}
	_ json.Marshaler     = Pod{}
)

func (pod Pod) RenderJSON() (string, error) {
	data, err := pod.MarshalJSON()
	return string(data), err
}

func (pod Pod) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(pod.origin, "", model.Indent)
	return data, err
}
