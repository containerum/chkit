package volume

import (
	"encoding/json"
	"strings"
)

var (
	_ json.Marshaler = new(Volume)
)
var (
	indent = strings.Repeat(" ", 4)
)

func (v *Volume) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(v, "", indent)
	return data, err
}
