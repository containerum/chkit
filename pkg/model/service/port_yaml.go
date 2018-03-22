package service

import (
	"github.com/containerum/chkit/pkg/model"
	yaml "gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = Port{}
)

func (port Port) RenderYAML() (string, error) {
	data, err := yaml.Marshal(port)
	return string(data), err
}
