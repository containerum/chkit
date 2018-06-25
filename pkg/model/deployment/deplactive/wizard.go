package deplactive

import (
	"fmt"
	"strings"

	"strconv"

	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/model/limits"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/validation"
)

type Wizard struct {
	Deployment *deployment.Deployment
	EditName   bool
	Configmaps []string
	Volumes    []string
}

func (config Wizard) Run() deployment.Deployment {
	if config.Deployment == nil {
		var depl = &deployment.Deployment{}
		Fill(depl)
		config.Deployment = depl
	}
	for exit := false; !exit; {
		var menuItems activekit.MenuItems
		if config.EditName {
			menuItems = activekit.MenuItems{componentEditNameMenu(config.Deployment)}
		}
		menuItems = menuItems.Append(componentEditReplicas(config.Deployment)).
			Append(componentEditContainers(config)...).
			Append(&activekit.MenuItem{
				Label: "Confirm",
				Action: func() error {
					if err := ValidateDeployment(*config.Deployment); err != nil {
						fmt.Println(err)
					} else {
						exit = true
					}
					return nil
				},
			})
		(&activekit.Menu{
			Title: fmt.Sprintf("Edit deployment"),
			Items: menuItems,
		}).Run()
	}
	return *config.Deployment
}

func componentEditNameMenu(deployment *deployment.Deployment) *activekit.MenuItem {
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Edit name : %s",
			activekit.OrString(deployment.Name, "undefined, required")),
		Action: func() error {
			for {
				var name = activekit.Promt("Type deployment name (hit Enter to leave %s): ",
					activekit.OrString(deployment.Name, "empty"))
				name = strings.TrimSpace(name)
				if err := validation.ValidateLabel(name); name != "" && err == nil {
					deployment.Name = name
				} else if name != "" && err != nil {
					fmt.Printf("%s is invalid deployment name\n", name)
					continue
				}
				break
			}
			return nil
		},
	}
}

func componentEditReplicas(deployment *deployment.Deployment) *activekit.MenuItem {
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Edit replicas : %d", deployment.Replicas),
		Action: func() error {
			for {
				var replicasStr = activekit.Promt("Type number of replicas to use (hit Enter to use %d, expected number in %v): ",
					deployment.Replicas, limits.ReplicasLimit)
				replicasStr = strings.TrimSpace(replicasStr)
				if replicas, err := strconv.Atoi(replicasStr); replicasStr != "" && err == nil {
					if !limits.ReplicasLimit.Containing(replicas) {
						fmt.Printf("Replicas number must be number in %v\n", limits.ReplicasLimit)
						continue
					}
					deployment.Replicas = replicas
				} else if err != nil {
					fmt.Printf("%q is invalid replicas number\n", replicasStr)
					continue
				}
				break
			}
			return nil
		},
	}
}
