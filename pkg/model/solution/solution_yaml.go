package solution

import (
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = UserSolution{}
	_ yaml.Marshaler     = UserSolution{}
)

func (solution UserSolution) RenderYAML() (string, error) {
	data, err := yaml.Marshal(solution)
	return string(data), err
}

func (solution UserSolution) MarshalYAML() (interface{}, error) {
	return solution.ToKube(), nil
}
