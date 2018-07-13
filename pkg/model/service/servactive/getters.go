package servactive

import (
	"fmt"

	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/validation"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/namegen"
)

const (
	ErrInvalidPort chkitErrors.Err = "invalid port"
)

func getName(defaultName string) string {
	for {
		name := activekit.Promt(fmt.Sprintf("Type service name (just leave empty to dub it %s)",
			defaultName))
		if name == "" {
			return defaultName
		}
		if err := validation.ValidateLabel(name); err != nil {
			fmt.Printf("\nError: %v\nPrint new one: ", err)
			continue
		}
		return name
	}

}

func editPorts(ports []service.Port, external bool) []service.Port {
	oldPorts := make([]service.Port, len(ports))
	copy(oldPorts, ports)
	var ok bool
	for exit := false; !exit; {
		var menu []*activekit.MenuItem
		for i, port := range ports {
			menu = append(menu, &activekit.MenuItem{
				Label: fmt.Sprintf("Edit port %q", port.Name),
				Action: func(i int) func() error {
					return func() error {
						port, deletePort := portEditorWizard(ports[i], external)
						if deletePort {
							ports = append(ports[:i], ports[i+1:]...)
						} else {
							ports[i] = port
						}
						return nil
					}
				}(i),
			})
		}
		(&activekit.Menu{
			Items: append(menu, []*activekit.MenuItem{
				{
					Label: "Add port",
					Action: func() error {
						ports = append(ports, portCreationWizard(service.Port{
							Name:       namegen.Aster() + "-" + namegen.Color(),
							TargetPort: 80,
							Protocol:   "TCP",
						}, external))
						return nil
					},
				},
				{
					Label: "Confirm",
					Action: func() error {
						ok = true
						exit = true
						return nil
					},
				},
				{
					Label: "Return to previous menu",
					Action: func() error {
						ok = false
						exit = true
						return nil
					},
				},
			}...),
		}).Run()
	}
	if ok {
		return ports
	}
	return oldPorts
}

func validatePort(port service.Port) error {
	var errs []error
	if err := validation.ValidateLabel(port.Name); port.Name == "" || (err != nil) {
		errs = append(errs, fmt.Errorf("\n + invalid port name %q", port.Name))
	}
	if port.Port != nil {
		if *port.Port < 1 || *port.Port > 65535 {
			errs = append(errs, fmt.Errorf("\n + invalid port %d: must be 1..65535", *port.Port))
		}
	}
	if port.TargetPort < 1 || port.TargetPort > 65535 {
		errs = append(errs, fmt.Errorf("\n + invalid target port %d: must be 1..65535", port.TargetPort))
	}
	if port.Protocol != "TCP" && port.Protocol != "UDP" {
		errs = append(errs, fmt.Errorf("\n + invalid port protocol: must be TCP or UDP"))
	}
	if len(errs) == 0 {
		return nil
	}
	return ErrInvalidPort.CommentF("name=%q", port.Name).AddReasons(errs...)
}

func getDeploy(defaultDepl string, depls []string) string {
	var menu []*activekit.MenuItem
	selectedDepl := defaultDepl
	for _, depl := range depls {
		menu = append(menu, &activekit.MenuItem{
			Label: depl,
			Action: func(depl string) func() error {
				return func() error {
					selectedDepl = depl
					return nil
				}
			}(depl),
		})
	}
	(&activekit.Menu{
		Items: append(menu, []*activekit.MenuItem{
			{
				Label: "Use custom deployment",
				Action: func() error {
					deployment := activekit.Promt("Type deployment label: ")
					if deployment == "" {
						return nil
					}
					if err := validation.ValidateLabel(deployment); err != nil {
						fmt.Printf("Invalid deployment label :(\n")
						return nil
					}
					selectedDepl = deployment
					return nil
				},
			},
			{
				Label: "Return to previous menu",
			},
		}...),
	}).Run()
	return selectedDepl
}
