package service

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = ServiceList{}
)

func (list ServiceList) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(list, "", model.Indent)
	return string(data), err
}
