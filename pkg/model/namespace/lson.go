package namespace

import "encoding/json"

var (
	_ json.Marshaler = new(NamespaceList)
)

func (list NamespaceList) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(list, "", indent)
	return data, err
}
