package pod

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = new(Pod)
	_ json.Marshaler     = new(JSONpod)
)

func (pod *Pod) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(pod, "", "   ")
	return string(data), err
}

type JSONpod struct {
	Pod
}

func (pod JSONpod) MarshalJSON() ([]byte, error) {
	data, err := pod.RenderJSON()
	return []byte(data), err
}
