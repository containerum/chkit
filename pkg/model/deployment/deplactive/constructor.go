package deplactive

import (
	"fmt"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/activeToolkit"
	"github.com/containerum/chkit/pkg/util/namegen"
)

const (
	ErrUserStopped chkitErrors.Err = "user stopped"
)

type Config struct {
	Force bool
}

func RunInteractveConstructor(config Config) (deployment.DeploymentList, error) {

	return nil, nil
}

func constructDeployment(config Config) (deployment.Deployment, error) {
	depl := defaultDeployment()
	for {
		_, n, _ := activeToolkit.Options("Whats't next?", false,
			fmt.Sprintf("Set name     : %s", depl.Name),
			fmt.Sprintf("Set replicas : %d", depl.Replicas),
			fmt.Sprintf("Set containers: %v", func() string {
				if len(depl.Containers) == 0 {
					return "none (required)"
				} else {
					return fmt.Sprintf("%v", depl.Containers)
				}
			}()),
			"Confirm",
			"Exit")
		switch n {
		case 0:
			depl.Name = getName(depl.Name)
		case 1:
			depl.Replicas = getReplicas(depl.Replicas)
		case 2:

		case 3:
			return depl, nil
		default:
			return depl, ErrUserStopped
		}
	}
}

func defaultDeployment() deployment.Deployment {
	return deployment.Deployment{
		Name:       namegen.Color() + "-" + namegen.Aster(),
		Replicas:   1,
		Containers: nil,
	}
}
