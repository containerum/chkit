package ingress

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
	model2 "github.com/containerum/kube-client/pkg/model"
)

var (
	_ model.JSONrenderer = Ingress{}
	_ json.Marshaler     = Ingress{}
	_ json.Unmarshaler   = new(Ingress)
)

func (ingr Ingress) RenderJSON() (string, error) {
	data, err := ingr.MarshalJSON()
	return string(data), err
}

func (ingr Ingress) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(ingr.ToKube(), "", "  ")
	return data, err
}

func (ingr *Ingress) UnmarshalJSON(data []byte) error {
	var kubeIngr model2.Ingress
	if err := json.Unmarshal(data, &kubeIngr); err != nil {
		return err
	}
	*ingr = IngressFromKube(kubeIngr)
	return nil
}
