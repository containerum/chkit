package namespace

import (
	"encoding/json"
	"strings"
)

var (
	_ json.Marshaler = new(Namespace)
)

var (
	indent = strings.Repeat(" ", 4)
)

func (ns *Namespace) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(ns, "", indent)
	return data, err
}
