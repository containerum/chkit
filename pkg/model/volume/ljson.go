package volume

import "encoding/json"

var (
	_ json.Marshaler = new(VolumeList)
)

func (list VolumeList) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(list, "", indent)
	return data, err
}
