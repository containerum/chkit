package solution

import (
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = SolutionEnv{}
)

func (envs SolutionEnv) RenderYAML() (string, error) {
	data, err := yaml.Marshal(envs)
	return string(data), err
}
