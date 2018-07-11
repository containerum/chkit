package container

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
	kubeModels "github.com/containerum/kube-client/pkg/model"
)

var (
	_ json.Marshaler   = Container{}
	_ json.Unmarshaler = new(Container)
)

func (cont Container) RenderJSON() (string, error) {
	var data, err = cont.MarshalJSON()
	return string(data), err
}

func (cont Container) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(cont.ToKube(), "", model.Indent)
}

func (cont *Container) UnmarshalJSON(data []byte) error {
	var kubeDepl kubeModels.Container
	if err := json.Unmarshal(data, &kubeDepl); err != nil {
		return err
	}
	*cont = Container{kubeDepl}
	return nil
}
