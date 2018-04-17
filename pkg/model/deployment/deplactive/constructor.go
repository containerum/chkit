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
	ErrInvalidContainer  chkitErrors.Err = "invalid container"
	ErrInvalidDeployment chkitErrors.Err = "invalid deployment"
)

type Config struct {
	Force      bool
	Deployment *deployment.Deployment
}

func Wizard(config Config) (deployment.Deployment, error) {
	var depl deployment.Deployment
	if config.Deployment == nil {
		depl = DefaultDeployment()
	} else {
		depl = *config.Deployment
	}
	for exit := false; !exit; {
		_, err := (&activekit.Menu{
			Items: []*activekit.MenuItem{
				{
					Label: fmt.Sprintf("Set name     : %s", depl.Name),
					Action: func() error {
						depl.Name = getName(depl.Name)
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set replicas : %d", depl.Replicas),
					Action: func() error {
						depl.Replicas = getReplicas(depl.Replicas)
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set containers: %v", activekit.OrValue(depl.Containers, "none (required)")),
					Action: func() error {
						depl.Containers = getContainers(depl.Containers)
						return nil
					},
				},
				{
					Label: "From file",
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
				{
					Label: "Confirm",
					Action: func() error {
						if err := validateDeployment(depl); err != nil {
							errTxt := err.Error()
							width := 0
							for _, line := range strings.Split(errTxt, "\n") {
								if len(line) > width {
									width = len(line)
								}
							}
							attention := strings.Repeat("!", width)
							fmt.Printf("%s\n%v\n%s\n", attention, err, attention)
							return nil
						}
						exit = true
						return nil
					},
				},
				{
					Label: "Exit",
					Action: func() error {
						if yes, _ := activekit.Yes("Are you sure you want to exit?"); yes {
							os.Exit(0)
						}
						return nil
					},
				},
			},
		}).Run()
		if err != nil {
			return depl, err
		}
	}
	return depl, nil
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

	if cont.Limits.CPU < 1 || cont.Limits.CPU > 12000 {
		errs = append(errs, fmt.Errorf("\n + invald CPU limit %d: must be in 0.001..12 milliCPU", cont.Limits.CPU))
	}

	if cont.Limits.Memory < 1 || cont.Limits.Memory > 16000 {
		errs = append(errs, fmt.Errorf("\n + invalid memory limit: must be in 1..16000 Mb"))
	}

	if len(errs) > 0 {
		return ErrInvalidContainer.CommentF("label=%q", cont.Name).AddReasons(errs...)
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
		return ErrInvalidDeployment.CommentF("label=%q", depl.Name).AddReasons(errs...)
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
