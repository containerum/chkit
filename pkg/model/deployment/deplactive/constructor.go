package deplactive

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/activeToolkit"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/chkit/pkg/util/validation"
)

const (
	ErrUserStopped       chkitErrors.Err = "user stopped"
	ErrInvalidDeployment chkitErrors.Err = "invalid deployment"
	ErrInvalidContainer  chkitErrors.Err = "invalid container"
)

type Config struct {
	Force      bool
	Deployment *deployment.Deployment
}

func ConstructDeployment(config Config) (deployment.Deployment, error) {
	var depl deployment.Deployment
	if config.Deployment == nil {
		depl = defaultDeployment()
	} else {
		depl = *config.Deployment
	}
	for {
		_, n, _ := activeToolkit.Options("Whats't next?", false,
			fmt.Sprintf("Set   name     : %s", depl.Name),
			fmt.Sprintf("Set   replicas : %d", depl.Replicas),
			fmt.Sprintf("Edit  containers: %v", activeToolkit.OrValue(depl.Containers, "none (required)")),
			"From file",
			"Confirm",
			"Exit")
		switch n {
		case 0:
			depl.Name = getName(depl.Name)
		case 1:
			depl.Replicas = getReplicas(depl.Replicas)
		case 2:
			depl.Containers = getContainers(depl.Containers)
		case 3:
			if filename, _ := activeToolkit.AskLine("print filename > "); strings.TrimSpace(filename) == "" {
				fmt.Printf("No file chosen\n")
				continue
			} else {
				var err error
				depl, err = FromFile(filename)
				if err != nil {
					return depl, err
				}
				continue
			}
		case 4:
			if err := validateDeployment(depl); err != nil {
				fmt.Printf("\n%v\n", err)
				continue
			}
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

func validateDeployment(depl deployment.Deployment) error {
	var errs []error
	if err := validation.ValidateLabel(depl.Name); err != nil {
		errs = append(errs, err)
	}
	if depl.Replicas < 1 || depl.Replicas > 15 {
		errs = append(errs, fmt.Errorf("invalid number of replicas %d: must be in 1..15", depl.Replicas))
	}
	for _, cont := range depl.Containers {
		conterr := validateContainer(cont)
		if conterr != nil {
			errs = append(errs, fmt.Errorf("invalid container %q: %v", cont.Name, conterr))
		}
	}
	return ErrInvalidDeployment.Wrap(errs...)
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
