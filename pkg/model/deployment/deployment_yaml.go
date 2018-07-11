package deployment

import (
	"github.com/containerum/chkit/pkg/model"
	model2 "github.com/containerum/kube-client/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = Deployment{}
	_ yaml.Marshaler     = Deployment{}
	_ yaml.Unmarshaler   = new(Deployment)
)

func (depl Deployment) RenderYAML() (string, error) {
	data, err := yaml.Marshal(depl)
	return string(data), err
}

func (depl Deployment) MarshalYAML() (interface{}, error) {
	return depl.ToKube(), nil
}

func (depl *Deployment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var kubeDepl model2.Deployment
	if err := unmarshal(&kubeDepl); err != nil {
		return err
	}
	*depl = DeploymentFromKube(kubeDepl)
	return nil
}
