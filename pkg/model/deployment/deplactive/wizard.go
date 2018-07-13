package deplactive

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/model/limits"
	"github.com/containerum/chkit/pkg/porta"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/containerum/chkit/pkg/util/validation"
	"github.com/sirupsen/logrus"
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
		menuItems = menuItems.Append(componentEditReplicas(config.Deployment),
			componentEditVersion(config.Deployment)).
			Append(componentEditContainers(config)...).
			Append(&activekit.MenuItem{
				Label: "Print to terminal",
				Action: func() error {
					data, err := config.Deployment.RenderYAML()
					if err != nil {
						logrus.WithError(err).Errorf("unable to render container to yaml")
						activekit.Attention(err.Error())
					}
					border := strings.Repeat("_", text.Width(data))
					fmt.Printf("%s\n%s\n%s\n", border, data, border)
					return nil
				},
			},
				&activekit.MenuItem{
					Label: "Export container to file",
					Action: func() error {
						var fname = activekit.Promt("Type filename: ")
						fname = strings.TrimSpace(fname)
						if fname != "" {
							if err := (porta.Exporter{OutFile: fname}.Export(config.Deployment)); err != nil {
								ferr.Printf("unable to export configmap:\n%v\n", err)
							}
						}
						return nil
					},
				},
				&activekit.MenuItem{
					Label: "Confirm",
					Action: func() error {
						if err := ValidateDeployment(*config.Deployment); err != nil {
							fmt.Println(err)
						} else {
							exit = true
						}
						return nil
					},
				},
				&activekit.MenuItem{
					Label: "Exit",
					Action: func() error {
						os.Exit(0)
						return nil
					},
				},
			)
		(&activekit.Menu{
			Title: fmt.Sprintf("Edit deployment"),
			Items: menuItems,
		}).Run()
	}
	return *config.Deployment
}

func componentEditVersion(deployment *deployment.Deployment) *activekit.MenuItem {
	var item = &activekit.MenuItem{
		Label: fmt.Sprintf("Edit version  : %s", deployment.Version),
		Action: func() error {
			(&activekit.Menu{
				Title: "Select version from container or set custom",
				Items: activekit.ItemsFromIter(uint(len(deployment.Containers)), func(index uint) *activekit.MenuItem {
					var cont = deployment.Containers[index].Copy()
					var version, ok = cont.SemanticVersion()
					if ok {
						return &activekit.MenuItem{
							Label: fmt.Sprintf("%v (%s)", cont.Version(), cont.Name),
							Action: func() error {
								deployment.Version = version
								return nil
							},
						}
					}
					return nil
				}).Append(
					&activekit.MenuItem{
						Label: "Set custom version",
						Action: func() error {
							for exit := false; !exit; {
								var vStr = activekit.Promt("Type version, v2.3.4, 1.4.2, etc., hit Enter to leave %s: ", deployment.Version)
								vStr = strings.TrimSpace(vStr)
								var version, err = semver.ParseTolerant(vStr)
								if err != nil {
									fmt.Printf("Invalid version string: %v", err)
									continue
								}
								deployment.Version = version
								return nil
							}
							return nil
						},
					},
					&activekit.MenuItem{
						Label: "Return to previous menu",
						Action: func() error {
							return nil
						},
					}),
			}).Run()
			return nil
		},
	}
	return item
}

func componentEditNameMenu(deployment *deployment.Deployment) *activekit.MenuItem {
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Edit name     : %s",
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
