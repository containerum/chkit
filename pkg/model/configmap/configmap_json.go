package configmap

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
	kubeModel "github.com/containerum/kube-client/pkg/model"
)

var (
	_ model.JSONrenderer = ConfigMap{}
	_ json.Marshaler     = ConfigMap{}
)

func (config ConfigMap) RenderJSON() (string, error) {
	data, err := config.MarshalJSON()
	return string(data), err
}

func (config ConfigMap) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(config.ToKube(), "", model.Indent)
}

func (config *ConfigMap) UnmarshalJSON(data []byte) error {
	var kubeCm kubeModel.ConfigMap
	if err := json.Unmarshal(data, &kubeCm); err != nil {
		return err
	}
	*config = ConfigMapFromKube(kubeCm)
	return nil
}
