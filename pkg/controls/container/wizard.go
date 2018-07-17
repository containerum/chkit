package container

import (
	"fmt"
	"os"
	"strings"

	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/porta"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/containerum/chkit/pkg/util/validation"
	"github.com/ninedraft/boxofstuff/str"
	"github.com/sirupsen/logrus"
)

type Wizard struct {
	Container   container.Container
	EditName    bool
	Deployment  string
	Deployments str.Vector
	Volumes     str.Vector
	Configs     str.Vector
}

func (wizard Wizard) Run() container.Container {
	var cont = wizard.Container.Copy()
	var items activekit.MenuItems
	if wizard.EditName {
		items = activekit.MenuItems{componentName(&cont)}
	}
	for exit := false; !exit; {
		(&activekit.Menu{
			Title: "Container " + cont.Name,
			Items: items.Append(
				componentDeployment(&cont, &wizard.Deployment, wizard.Deployments),
				componentImage(&cont),
				componentLimits(&cont),
				componentVolumes(&cont, wizard.Volumes),
				componentConfigmaps(&cont, wizard.Configs),
				componentEnvs(&cont),
				componentCmds(&cont),
				&activekit.MenuItem{
					Label: "Print to terminal",
					Action: func() error {
						data, err := cont.RenderYAML()
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
							if err := (porta.Exporter{OutFile: fname}.Export(cont)); err != nil {
								ferr.Printf("unable to export configmap:\n%v\n", err)
							}
						}
						return nil
					},
				},
				&activekit.MenuItem{
					Label: "Confirm",
					Action: func() error {
						if err := cont.Validate(); err != nil {
							fmt.Println(err)
							return nil
						}
						exit = true
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
			),
		}).Run()
	}
	return cont
}

func componentName(cont *container.Container) *activekit.MenuItem {
	return &activekit.MenuItem{
		Label: "Container name : " + str.Vector{cont.Name, "undefined, required"}.FirstNonEmpty(),
		Action: func() error {
			for {
				var name = activekit.Promt("Type container name, hit Enter to leave %s: ", str.Vector{cont.Name,
					"empty"}.FirstNonEmpty())
				name = strings.TrimSpace(name)
				if err := validation.ValidateLabel(name); name != "" && err == nil {
					cont.Name = name
				} else if name != "" && err != nil {
					fmt.Printf("invalid container name: %v", err)
					continue
				}
				return nil
			}
		},
	}
}

func componentDeployment(cont *container.Container, depl *string, deployments str.Vector) *activekit.MenuItem {
	var oldDepl = *depl
	return &activekit.MenuItem{
		Label: "Deployment     : " + str.Vector{*depl, "undefined, required"}.FirstNonEmpty(),
		Action: func() error {
			(&activekit.Menu{
				Title: "Container -> Select deployment",
				Items: activekit.StringSelector(deployments, func(s string) error {
					*depl = s
					return nil
				}).Append(&activekit.MenuItem{
					Label: "Confirm",
					Action: func() error {
						*depl = oldDepl
						return nil
					},
				},
					&activekit.MenuItem{
						Label: "Return to previous menu, use deployment" + oldDepl,
						Action: func() error {
							return nil
						},
					}),
			}).Run()
			return nil
		},
	}
}

func componentImage(cont *container.Container) *activekit.MenuItem {
	var item = &activekit.MenuItem{
		Label: "Edit image     : " + str.Vector{cont.Image, "undefined, required"}.FirstNonEmpty(),
		Action: func() error {
			for {
				var image = activekit.Promt("Type container image, hit Enter to leave %s: ", str.Vector{cont.Image,
					"empty"}.FirstNonEmpty())
				image = strings.TrimSpace(image)
				if err := validation.ValidateImageName(image); image != "" && err == nil {
					cont.Image = image
				} else if image != "" && err != nil {
					fmt.Printf("invalid container image: %v", err)
					continue
				}
				return nil
			}
		},
	}
	return item
}
