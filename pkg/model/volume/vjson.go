package volume

import (
	"encoding/json"
	"strings"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.JSONrenderer = new(Volume)
)
var (
	indent = strings.Repeat(" ", 4)
)

func (v *Volume) RenderJSON() (string, error) {
	data, err := json.MarshalIndent(v, "", indent)
	return string(data), err
}
