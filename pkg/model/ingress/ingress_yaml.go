package ingress

import (
	"github.com/containerum/chkit/pkg/model"
	model2 "github.com/containerum/kube-client/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = Ingress{}
	_ yaml.Marshaler     = Ingress{}
	_ yaml.Unmarshaler   = new(Ingress)
)

func (ingr Ingress) RenderYAML() (string, error) {
	data, err := yaml.Marshal(ingr)
	return string(data), err
}

func (ingr Ingress) MarshalYAML() (interface{}, error) {
	return ingr.ToKube(), nil
}

func (ingr *Ingress) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var kubeIngr model2.Ingress
	if err := unmarshal(&kubeIngr); err != nil {
		return err
	}
	*ingr = IngressFromKube(kubeIngr)
	return nil
}
