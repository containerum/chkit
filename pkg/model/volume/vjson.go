package volume

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = new(Volume)
	_ json.Marshaler     = Volume{}
)

func (vol Volume) RenderJSON() (string, error) {
	data, err := vol.MarshalJSON()
	return string(data), err
}

func (vol Volume) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(vol.origin, "", "    ")
	return data, err
}
