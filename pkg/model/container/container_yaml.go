package container

import (
	kubeModels "github.com/containerum/kube-client/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ yaml.Marshaler   = Container{}
	_ yaml.Unmarshaler = new(Container)
)

func (container Container) RenderYAML() (string, error) {
	var data, err = yaml.Marshal(container)
	return string(data), err
}

func (container Container) MarshalYAML() (interface{}, error) {
	return container.ToKube(), nil
}

func (container *Container) UnmarshalYAML(decode func(interface{}) error) error {
	var cont kubeModels.Container
	if err := decode(&cont); err != nil {
		return err
	}
	*container = Container{cont}
	return nil
}
