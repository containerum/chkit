package deplactive

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/containerum/chkit/pkg/model/deployment"
	"gopkg.in/yaml.v2"
)

func FromFile(filename string) (deployment.Deployment, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return deployment.Deployment{}, err
	}
	kubeDepl := (&deployment.Deployment{}).ToKube()
	switch filepath.Ext(filename) {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &kubeDepl)
	default:
		err = json.Unmarshal(data, &kubeDepl)
	}
	if err != nil {
		return deployment.Deployment{}, err
	}
	depl := deployment.DeploymentFromKube(kubeDepl)
	return depl, ValidateDeployment(depl)
}
