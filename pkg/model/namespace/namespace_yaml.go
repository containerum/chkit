package namespace

import (
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = Namespace{}
	_ yaml.Marshaler     = Namespace{}
)

func (ns Namespace) RenderYAML() (string, error) {
	data, err := yaml.Marshal(ns)
	return string(data), err
}

func (ns Namespace) MarshalYAML() (interface{}, error) {
	return ns.origin, nil
}
