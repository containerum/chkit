package configmap

import (
	"github.com/containerum/chkit/pkg/model"
	kubeModel "github.com/containerum/kube-client/pkg/model"
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

func (config *ConfigMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var kubeCm kubeModel.ConfigMap
	if err := unmarshal(&kubeCm); err != nil {
		return err
	}
	*config = ConfigMapFromKube(kubeCm)
	return nil
}
