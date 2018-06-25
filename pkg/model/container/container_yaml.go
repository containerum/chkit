package container

import "gopkg.in/yaml.v2"

var (
	_ yaml.Marshaler = Container{}
)

func (container Container) RenderYAML() (string, error) {
	var data, err = yaml.Marshal(container)
	return string(data), err
}

func (container Container) MarshalYAML() (interface{}, error) {
	return container.ToKube(), nil
}
