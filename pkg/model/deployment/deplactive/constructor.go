package deplactive

import (
	"fmt"
	"strings"

	"os"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/chkit/pkg/util/validation"
)

const (
	ErrUserStopped       chkitErrors.Err = "user stopped"
	ErrInvalidContainer  chkitErrors.Err = "invalid container"
	ErrInvalidDeployment chkitErrors.Err = "invalid deployment"
)

type Config struct {
	Force      bool
	Deployment *deployment.Deployment
}

func ConstructDeployment(config Config) (deployment.Deployment, error) {
	var depl deployment.Deployment
	if config.Deployment == nil {
		depl = DefaultDeployment()
	} else {
		depl = *config.Deployment
	}
	exit := &activekit.MenuItem{
		Name: "Exit",
	}
	confirm := &activekit.MenuItem{
		Name: "Confirm",
		Action: func() error {
			return validateDeployment(depl)
		},
	}
	for {
		result, err := (&activekit.Menu{
			Items: []*activekit.MenuItem{
				{
					Name: fmt.Sprintf("Set name     : %s", depl.Name),
					Action: func() error {
						depl.Name = getName(depl.Name)
						return nil
					},
				},
				{
					Name: fmt.Sprintf("Set replicas : %d", depl.Replicas),
					Action: func() error {
						depl.Replicas = getReplicas(depl.Replicas)
						return nil
					},
				},
				{
					Name: fmt.Sprintf("Set containers: %v", activekit.OrValue(depl.Containers, "none (required)")),
					Action: func() error {
						depl.Containers = getContainers(depl.Containers)
						return nil
					},
				},
				{
					Name: "From file",
					Action: func() error {
						if filename, _ := activekit.AskLine("print filename > "); strings.TrimSpace(filename) == "" {
							fmt.Printf("No file chosen\n")
							return nil
						} else {
							var err error
							depl, err = FromFile(filename)
							if err != nil {
								fmt.Println(err)
								return nil
							}
						}
						return nil
					},
				},
				confirm,
				exit,
			},
		}).Run()
		switch result {
		case exit:
			os.Exit(0)
		case confirm:
			if err != nil {
				errTxt := err.Error()
				width := 0
				for _, line := range strings.Split(errTxt, "\n") {
					if len(line) > width {
						width = len(line)
					}
				}
				attention := strings.Repeat("!", width)
				fmt.Printf("%s\n%v\n%s\n", attention, err, attention)
				continue
			}
			return depl, nil
		default:
			if err != nil {
				return depl, err
			}
		}
	}
}

func DefaultDeployment() deployment.Deployment {
	return deployment.Deployment{
		Name:       namegen.Color() + "-" + namegen.Aster(),
		Replicas:   1,
		Containers: nil,
	}
}

func validateContainer(cont container.Container) error {
	var errs []error
	if err := validation.ValidateLabel(cont.Name); err != nil {
		errs = append(errs, fmt.Errorf("\n + invalid container name: %v", err))
	}
	if err := validation.ValidateImageName(cont.Image); err != nil {
		errs = append(errs, fmt.Errorf("\n + invalid image name: %v", err))
	}
	if cont.Limits.CPU == "" {
		errs = append(errs, fmt.Errorf("\n + undefined CPU limit"))
	}
	if cont.Limits.Memory == "" {
		errs = append(errs, fmt.Errorf("\n + undefined memory limit"))
	}
	if len(errs) > 0 {
		return ErrInvalidContainer.CommentF("%q", cont.Name).AddReasons(errs...)
	}
	return nil
}

func validateDeployment(depl deployment.Deployment) error {
	var errs []error
	if depl.Replicas < 1 || depl.Replicas > 15 {
		errs = append(errs, fmt.Errorf("\n + invalid replicas number %d: must be 1..15", depl.Replicas))
	}
	if len(depl.Containers) == 0 {
		errs = append(errs, fmt.Errorf("\n + can't create deployment without containers!"))
	}
	for _, cont := range depl.Containers {
		if err := validateContainer(cont); err != nil {
			errs = append(errs, fmt.Errorf("\n + %s", indent("  ", err.Error())))
		}
	}
	if len(errs) > 0 {
		return ErrInvalidDeployment.Comment("\n").AddReasons(errs...)
	}
	return nil
}

func indent(indent, str string) string {
	var lines []string
	for i, line := range strings.Split(str, "\n") {
		if i != 0 {
			line = indent + line
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return indent + str
	}
	return strings.Join(lines, "\n")
}
