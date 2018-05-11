package configmap

import (
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = ConfigMap{}
	_ yaml.Marshaler     = ConfigMap{}
)

func (config ConfigMap) RenderYAML() (string, error) {
	data, err := yaml.Marshal(config)
	return string(data), err
}

func (config ConfigMap) MarshalYAML() (interface{}, error) {
	return config.ToKube(), nil
}
