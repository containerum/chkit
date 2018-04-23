package deplactive

import (
	"fmt"
	"os"
	"strings"

	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/activekit"
)

func ReplaceWizard(config Config) (deployment.Deployment, error) {
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
						if err := ValidateDeployment(depl); err != nil {
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
