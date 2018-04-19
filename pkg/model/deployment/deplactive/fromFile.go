package deplactive

import (
	"encoding/json"
	"io/ioutil"

	"github.com/containerum/chkit/pkg/model/deployment"
)

func FromFile(filename string) (deployment.Deployment, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return deployment.Deployment{}, err
	}
	kubeDepl := (&deployment.Deployment{}).ToKube()
	err = json.Unmarshal(data, &kubeDepl)
	if err != nil {
		return deployment.Deployment{}, err
	}
	depl := deployment.DeploymentFromKube(kubeDepl)
	return depl, ValidateDeployment(depl)
}
