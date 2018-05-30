package namespace

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = new(NamespaceList)
)

func (list NamespaceList) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(list.ToKube(), "", model.Indent)
	return string(data), err
}
