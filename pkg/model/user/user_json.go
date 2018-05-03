package user

import (
	"github.com/containerum/chkit/pkg/model"
	"github.com/gin-gonic/gin/json"
)

var (
	_ model.JSONrenderer = User{}
)

func (user User) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(user.ToKube(), "", model.Indent)
	return string(data), err
}
