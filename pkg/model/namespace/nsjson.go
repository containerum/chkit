package namespace

import (
	"encoding/json"
	"strings"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = new(Namespace)
)

var (
	indent = strings.Repeat(" ", 4)
)

func (ns *Namespace) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(ns, "", indent)
	return string(data), err
}
