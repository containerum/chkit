package deplactive

import (
	"encoding/json"
	"io/ioutil"

	"github.com/containerum/chkit/pkg/model/deployment"
)

func fromFile(filename string) (deployment.Deployment, error) {
	depl := deployment.Deployment{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return depl, err
	}
	kubeDepl := depl.ToKube()
	err = json.Unmarshal(data, &kubeDepl)
	if err != nil {
		return depl, err
	}
	depl = deployment.DeploymentFromKube(kubeDepl)
	return depl, nil
}
