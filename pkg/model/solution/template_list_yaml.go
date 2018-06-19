package solution

import (
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = SolutionList{}
	_ yaml.Marshaler     = SolutionList{}
)

func (list SolutionList) RenderYAML() (string, error) {
	data, err := yaml.Marshal(list)
	return string(data), err
}

func (list SolutionList) MarshalYAML() (interface{}, error) {
	return list, nil
}
