package pod

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = new(PodList)
)

func (list PodList) RenderJSON() (string, error) {
	jsonList := make([]JSONpod, 0, len(list))
	for _, pod := range list {
		jsonList = append(jsonList, JSONpod{pod})
	}
	data, err := json.MarshalIndent(jsonList, "", "    ")
	return string(data), err
}
