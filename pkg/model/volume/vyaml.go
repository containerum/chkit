package volume

import (
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = Volume{}
	_ yaml.Marshaler     = Volume{}
)

func (v Volume) RenderYAML() (string, error) {
	data, err := yaml.Marshal(v)
	return string(data), err
}

func (vol Volume) MarshalYAML() (interface{}, error) {
	return vol.origin, nil
}
