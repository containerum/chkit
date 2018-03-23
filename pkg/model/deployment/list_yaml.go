package deployment

import (
	"github.com/containerum/chkit/pkg/model"
	yaml "gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = DeploymentList{}
)

func (list DeploymentList) RenderYAML() (string, error) {
	data, err := yaml.Marshal(list)
	return string(data), err
}
