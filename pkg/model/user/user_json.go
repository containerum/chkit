package user

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = User{}
)

func (user User) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(user.ToKube(), "", model.Indent)
	return string(data), err
}
