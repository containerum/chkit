package service

import (
	"github.com/containerum/chkit/pkg/model"
	kubeModel "github.com/containerum/kube-client/pkg/model"
	"gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = Service{}
	_ yaml.Marshaler     = Service{}
)

func (serv Service) RenderYAML() (string, error) {
	serv.ToKube()
	data, err := yaml.Marshal(serv)
	return string(data), err
}

func (serv Service) MarshalYAML() (interface{}, error) {
	return serv.ToKube(), nil
}

func (serv *Service) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var kubeSvc kubeModel.Service
	if err := unmarshal(&kubeSvc); err != nil {
		return err
	}
	*serv = ServiceFromKube(kubeSvc)
	return nil
}
