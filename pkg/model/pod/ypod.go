package pod

import (
	"github.com/containerum/chkit/pkg/model"
	yaml "gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = new(Pod)
)

func (pod *Pod) RenderYAML() (string, error) {
	data, err := yaml.Marshal(pod)
	return string(data), err
}
