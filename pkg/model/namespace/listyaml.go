package namespace

import (
	"github.com/containerum/chkit/pkg/model"
	yaml "gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = new(NamespaceList)
)

func (list NamespaceList) RenderYAML() (string, error) {
	data, err := yaml.Marshal(list)
	return string(data), err
}
