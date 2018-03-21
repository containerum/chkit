package volume

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = new(VolumeList)
)

func (list VolumeList) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(list, "", indent)
	return string(data), err
}
