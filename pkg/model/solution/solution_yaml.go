package solution

import (
	"github.com/containerum/chkit/pkg/model"
	kubeModels "github.com/containerum/kube-client/pkg/model"
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

func (solution *Solution) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var kubeSol kubeModels.UserSolution
	if err := unmarshal(&kubeSol); err != nil {
		return err
	}
	*solution = SolutionFromKube(kubeSol)
	return nil
}
