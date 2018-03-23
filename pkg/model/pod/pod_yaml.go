package pod

import (
	"github.com/containerum/chkit/pkg/model"
	yaml "gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = Pod{}
	_ yaml.Marshaler     = Pod{}
)

func (pod Pod) RenderYAML() (string, error) {
	data, err := yaml.Marshal(pod)
	return string(data), err
}

func (pod Pod) MarshalYAML() (interface{}, error) {
	return pod.origin, nil
}
