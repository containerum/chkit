package namespace

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = Namespace{}
	_ json.Marshaler     = Namespace{}
)

func (ns Namespace) RenderJSON() (string, error) {
	data, err := ns.MarshalJSON()
	return string(data), err
}

func (ns Namespace) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(ns.ToKube(), "", model.Indent)
	return data, err
}
