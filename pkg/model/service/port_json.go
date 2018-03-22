package service

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = Port{}
)

func (port Port) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(port, "", "    ")
	return string(data), err
}
