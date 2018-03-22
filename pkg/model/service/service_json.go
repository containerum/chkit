package service

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = Service{}
)

func (serv Service) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(serv, "", "    ")
	return string(data), err
}
