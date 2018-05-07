package solution

import (
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = Solution{}
	_ yaml.Marshaler     = Solution{}
)

func (solution Solution) RenderYAML() (string, error) {
	data, err := yaml.Marshal(solution)
	return string(data), err
}

func (solution Solution) MarshalYAML() (interface{}, error) {
	return solution.ToKube(), nil
}
