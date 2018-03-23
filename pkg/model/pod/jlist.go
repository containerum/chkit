package pod

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = new(PodList)
)

func (list PodList) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(list, "", "    ")
	return string(data), err
}
