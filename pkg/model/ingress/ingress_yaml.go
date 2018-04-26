package ingress

import (
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = Ingress{}
	_ yaml.Marshaler     = Ingress{}
)

func (ingr Ingress) RenderYAML() (string, error) {
	data, err := yaml.Marshal(ingr)
	return string(data), err
}

func (ingr Ingress) MarshalYAML() (interface{}, error) {
	return ingr.ToKube(), nil
}
