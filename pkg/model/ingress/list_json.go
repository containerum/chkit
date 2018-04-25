package ingress

import (
	"github.com/containerum/chkit/pkg/model"
	"github.com/gin-gonic/gin/json"
)

var (
	_ model.JSONrenderer = IngressList{}
)

func (list IngressList) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(list, "", "  ")
	return string(data), err
}
