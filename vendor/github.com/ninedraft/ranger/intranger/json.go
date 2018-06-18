package intranger

import "encoding/json"

type intRangerJSON struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

func (ranger *intRanger) UnmarshalJSON(p []byte) error {
	var jsonRanger intRangerJSON
	err := json.Unmarshal(p, &jsonRanger)
	if err != nil {
		return err
	}
	*ranger = IntRanger(jsonRanger.Min, jsonRanger.Max)
	return nil
}


func (ranger intRanger) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(ranger.toJSON())
	return data, err
}


func (ranger intRanger) toJSON() intRangerJSON {
	return intRangerJSON{
		Min: ranger.min,
		Max: ranger.max,
	}
}
