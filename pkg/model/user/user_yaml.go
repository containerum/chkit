package user

import (
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = User{}
)

func (user User) RenderYAML() (string, error) {
	data, err := yaml.Marshal(user.ToKube())
	return string(data), err
}
