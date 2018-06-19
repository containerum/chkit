package solution

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = TemplatesList{}
)

func (list TemplatesList) RenderJSON() (string, error) {
	data, err := list.MarshalJSON()
	return string(data), err
}

func (list TemplatesList) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(list, "", model.Indent)
	return data, err
}
