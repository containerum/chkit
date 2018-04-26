package ingress

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = Ingress{}
	_ json.Marshaler     = Ingress{}
)

func (ingr Ingress) RenderJSON() (string, error) {
	data, err := ingr.MarshalJSON()
	return string(data), err
}

func (ingr Ingress) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(ingr.ToKube(), "", "  ")
	return data, err
}
