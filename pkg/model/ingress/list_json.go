package ingress

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = IngressList{}
)

func (list IngressList) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(list, "", "  ")
	return string(data), err
}
