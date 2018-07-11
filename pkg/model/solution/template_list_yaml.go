package solution

import (
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = TemplatesList{}
)

func (list TemplatesList) RenderYAML() (string, error) {
	data, err := yaml.Marshal(list)
	return string(data), err
}
