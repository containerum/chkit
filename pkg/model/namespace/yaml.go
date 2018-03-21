package namespace

import (
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = new(Namespace)
)

func (ns *Namespace) RenderYAML() (string, error) {
	data, err := yaml.Marshal(ns)
	return string(data), err
}
