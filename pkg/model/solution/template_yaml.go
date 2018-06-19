package solution

import (
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = SolutionTemplate{}
	_ yaml.Marshaler     = SolutionTemplate{}
)

func (solution SolutionTemplate) RenderYAML() (string, error) {
	data, err := yaml.Marshal(solution)
	return string(data), err
}

func (solution SolutionTemplate) MarshalYAML() (interface{}, error) {
	return solution.ToKube(), nil
}
