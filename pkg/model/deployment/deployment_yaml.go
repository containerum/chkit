package deployment

import (
	"github.com/containerum/chkit/pkg/model"
	yaml "gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = Deployment{}
	_ yaml.Marshaler     = Deployment{}
)

func (depl Deployment) RenderYAML() (string, error) {
	data, err := yaml.Marshal(depl.ToKube())
	return string(data), err
}

func (depl Deployment) MarshalYAML() (interface{}, error) {
	return depl.origin, nil
}
