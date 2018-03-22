package deployment

import (
	"github.com/containerum/chkit/pkg/model"
	yaml "gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = Deployment{}
)

func (depl Deployment) RenderYAML() (string, error) {
	data, err := yaml.Marshal(depl)
	return string(data), err
}
