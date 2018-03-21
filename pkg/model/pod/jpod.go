package pod

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = new(Pod)
)

func (pod *Pod) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(pod, "", "   ")
	return string(data), err
}
