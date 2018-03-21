package volume

import (
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = new(Volume)
)

func (v *Volume) RenderYAML() (string, error) {
	data, err := yaml.Marshal(v)
	return string(data), err
}
