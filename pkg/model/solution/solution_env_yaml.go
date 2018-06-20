package solution

import (
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = SolutionEnv{}
	_ yaml.Marshaler     = Solution{}
)

func (solution SolutionEnv) RenderYAML() (string, error) {
	data, err := yaml.Marshal(solution)
	return string(data), err
}

func (solution SolutionEnv) MarshalYAML() (interface{}, error) {
	return solution, nil
}
