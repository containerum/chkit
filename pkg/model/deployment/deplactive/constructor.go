package deplactive

import (
	"fmt"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/activeToolkit"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/chkit/pkg/util/validation"
)

const (
	ErrUserStopped      chkitErrors.Err = "user stopped"
	ErrInvalidContainer chkitErrors.Err = "invalid container"
)

type Config struct {
	Force bool
}

func RunInteractveConstructor(config Config) (deployment.DeploymentList, error) {

	return nil, nil
}

func ConstructDeployment(config Config) (deployment.Deployment, error) {
	depl := defaultDeployment()
	for {
		_, n, _ := activeToolkit.Options("Whats't next?", false,
			fmt.Sprintf("Set name     : %s", depl.Name),
			fmt.Sprintf("Set replicas : %d", depl.Replicas),
			fmt.Sprintf("Set containers: %v", activeToolkit.OrValue(depl.Containers, "none (required)")),
			"Confirm",
			"Exit")
		switch n {
		case 0:
			depl.Name = getName(depl.Name)
		case 1:
			depl.Replicas = getReplicas(depl.Replicas)
		case 2:
			depl.Containers = getContainers()
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

func validateContainer(cont container.Container) error {
	var errs []error
	if err := validation.ValidateLabel(cont.Name); err != nil {
		errs = append(errs, err)
	}
	if err := validation.ValidateImageName(cont.Image); err != nil {
		errs = append(errs, err)
	}
	if cont.Limits.CPU == "" {
		errs = append(errs, fmt.Errorf("undefined CPU limit"))
	}
	if cont.Limits.Memory == "" {
		errs = append(errs, fmt.Errorf("undefined memory limit"))
	}
	if len(errs) > 0 {
		return ErrInvalidContainer.Wrap(errs...)
	}
	return nil
}
