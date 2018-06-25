package container

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
)

func (cont Container) RenderJSON() (string, error) {
	var data, err = cont.MarshalJSON()
	return string(data), err
}

func (cont Container) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(cont.ToKube(), "", model.Indent)
}
