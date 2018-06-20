package solution

import (
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = TemplatesList{}
	_ yaml.Marshaler     = TemplatesList{}
)

func (list TemplatesList) RenderYAML() (string, error) {
	data, err := yaml.Marshal(list)
	return string(data), err
}

func (list TemplatesList) MarshalYAML() (interface{}, error) {
	return list, nil
}
