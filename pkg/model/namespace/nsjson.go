package namespace

import (
	"encoding/json"
	"strings"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = Namespace{}
	_ json.Marshaler     = Namespace{}
)

var (
	indent = strings.Repeat(" ", 4)
)

func (ns Namespace) RenderJSON() (string, error) {
	data, err := ns.MarshalJSON()
	return string(data), err
}

func (ns Namespace) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(ns.origin, "", "    ")
	return data, err
}
