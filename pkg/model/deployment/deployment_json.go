package deployment

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
	model2 "github.com/containerum/kube-client/pkg/model"
)

var (
	_ model.JSONrenderer = Deployment{}
	_ json.Marshaler     = Deployment{}
	_ json.Unmarshaler   = new(Deployment)
)

func (depl Deployment) RenderJSON() (string, error) {
	depl.ToKube()
	data, err := depl.MarshalJSON()
	return string(data), err
}

func (depl Deployment) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(depl.ToKube(), "", model.Indent)
	return data, err
}

func (depl *Deployment) UnmarshalJSON(data []byte) error {
	var kubeDepl model2.Deployment
	if err := json.Unmarshal(data, &kubeDepl); err != nil {
		return err
	}
	*depl = DeploymentFromKube(kubeDepl)
	return nil
}
