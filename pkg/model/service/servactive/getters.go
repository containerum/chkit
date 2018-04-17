package servactive

import (
	"fmt"
	"net"
	"strings"

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

func getIPs(ips []string) []string {
	oldIPs := make([]string, len(ips))
	copy(oldIPs, ips)
	var ok bool
	for exit := false; !exit; {
		var menu []*activekit.MenuItem
		for i, ip := range ips {
			menu = append(menu, &activekit.MenuItem{
				Name: fmt.Sprintf("Delete %q", ip),
				Action: func(i int) func() error {
					return func() error {
						ips = append(ips[:i], ips[i+1:]...)
						return nil
					}
				}(i),
			})
		}
		menu = append(menu, []*activekit.MenuItem{
			{
				Name: "Add addr",
				Action: func() error {
					rawAddr := strings.TrimSpace(activekit.Promt("Type IP : "))
					ip := net.ParseIP(rawAddr)
					if ip == nil {
						fmt.Printf("\nInvalid IP address!\n")
						return nil
					}
					ips = append(ips, ip.String())
					return nil
				},
			},
			{
				Name: "Confirm",
				Action: func() error {
					exit = true
					ok = true
					return nil
				},
			},
			{
				Name: "Return to previous menu",
				Action: func() error {
					ok = false
					exit = true
					return nil
				},
			},
		}...)
		(&activekit.Menu{
			Items: menu,
		}).Run()
	}
	if ok {
		return ips
	}
	return oldIPs
}

func editPorts(ports []service.Port) []service.Port {
	oldPorts := make([]service.Port, len(ports))
	copy(oldPorts, ports)
	ok := false
	for exit := false; !exit; {
		var menu []*activekit.MenuItem
		for i, port := range ports {
			menu = append(menu, &activekit.MenuItem{
				Name: fmt.Sprintf("Edit port %q", port.Name),
				Action: func(i int) func() error {
					return func() error {
						port, deletePort := portEditorWizard(ports[i])
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
					Name: "Add port",
					Action: func() error {
						ports = append(ports, portCreationWizard(service.Port{
							Name:       namegen.Aster() + "-" + namegen.Color(),
							TargetPort: 80,
							Protocol:   "TCP",
						}))
						return nil
					},
				},
				{
					Name: "Confirm",
					Action: func() error {
						ok = true
						exit = true
						return nil
					},
				},
				{
					Name: "Return to previous menu",
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

func getPort(ports *[]service.Port, ind int) (service.Port, bool) {
	var p service.Port
	if ind < 0 || len(*ports) == 0 {
		p = service.Port{
			Name:     namegen.Aster() + "-" + namegen.Color(),
			Protocol: "TCP",
		}
		*ports = append(*ports, p)
	} else {
		p = (*ports)[ind]
	}
	ok := false
	for exit := false; !exit; {
		(&activekit.Menu{
			Items: []*activekit.MenuItem{
				{
					Name: fmt.Sprintf("Set name : %s",
						activekit.OrString(p.Name, "undefined (required)")),
					Action: func() error {
						var promt string
						if p.Name == "" {
							promt = "Print port name: "
						} else {
							promt = fmt.Sprintf("Print port name (hit enter to use %q): ", p.Name)
						}
						name := strings.TrimSpace(activekit.Promt(promt))
						if name == "" {
							return nil
						}
						if err := validation.ValidateLabel(name); err != nil {
							activekit.Attention(fmt.Sprintf("Invalid port name:\n%v", err))
							return nil
						}
						p.Name = name
						return nil
					},
				},
				{
					Name: fmt.Sprintf("Set target port : %d (required)", p.TargetPort),
					Action: func() error {
						promt := fmt.Sprintf("Print target port (1..65535, hit enter to use %d): ", p.TargetPort)
						portStr := strings.TrimSpace(activekit.Promt(promt))
						if portStr == "" {
							return nil
						}
						var port int
						if _, err := fmt.Sscan(portStr, "%d", &port); err != nil || (port < 1 && port > 65535) {
							fmt.Printf("Invalid target port %q: must be number 1..65535\n", portStr)
							return nil
						}
						p.TargetPort = port
						return nil
					},
				},
				{
					Name: fmt.Sprintf("Set proto : %s",
						activekit.OrString(p.Protocol, "undefined (required)")),
					Action: func() error {
						_, err := (&activekit.Menu{
							Title: fmt.Sprintf("Select protocol (current: %s)", p.Protocol),
							Items: []*activekit.MenuItem{
								{
									Name: "TCP",
									Action: func() error {
										p.Protocol = "TCP"
										return nil
									},
								},
								{
									Name: "UDP",
									Action: func() error {
										p.Protocol = "UDP"
										return nil
									},
								},
								{
									Name: "Return to previous menu",
								},
							}}).Run()
						return err
					},
				},
				{
					Name: fmt.Sprintf("Set port : %s",
						activekit.OrValue(p.Port, "undefined (optional)")),
					Action: func() error {
						var promt string
						if p.Port == nil {
							promt = "Print port (11000..65535, hit enter to leave empty):"
						} else {
							promt = fmt.Sprintf("Print port (11000..65535, hit enter to use %d, type 'none' to leave empty): ", *p.Port)
						}
						portStr := strings.TrimSpace(activekit.Promt(promt))
						if portStr == "none" {
							p.Port = nil
						}
						if portStr == "" {
							return nil
						}
						var port int
						if _, err := fmt.Sscan(portStr, "%d", &port); err != nil || (port < 11000 && port > 65535) {
							fmt.Printf("Invalid port %q: must be number in 11000..65535\n", portStr)
							return nil
						}
						p.Port = &port
						return nil
					},
				},
				{
					Name: fmt.Sprintf("Delete port %q", p),
					Action: func() error {
						*ports = append((*ports)[:ind], (*ports)[ind+1:]...)
						return nil
					},
				},
				{
					Name: "Confirm",
					Action: func() error {
						if err := validatePort(p); err != nil {
							activekit.Attention(err.Error())
							return nil
						}
						ok = true
						exit = false
						return nil
					},
				},
			},
		}).Run()
	}
	return p, ok
}

func validatePort(port service.Port) error {
	var errs []error
	if err := validation.ValidateLabel(port.Name); port.Name == "" || (err != nil) {
		errs = append(errs, fmt.Errorf("\n + invalid port name %q", port.Name))
	}
	if port.Port != nil {
		if *port.Port < 11000 || *port.Port > 65535 {
			errs = append(errs, fmt.Errorf("\n + invalid port %d: must be 11000..65535", *port.Port))
		}
	}
	if port.TargetPort < 1 || port.TargetPort > 65535 {
		errs = append(errs, fmt.Errorf("\n + invalid target port %d: must be 1..65535"))
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
			Name: depl,
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
				Name: "Use custom deployment",
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
				Name: "Return to previous menu",
			},
		}...),
	}).Run()
	return selectedDepl
}

func getDomain(defaultDomain string) string {
	domain := defaultDomain
	(&activekit.Menu{Items: []*activekit.MenuItem{
		{
			Name: "Use custom domain",
			Action: func() error {
				domain = activekit.Promt("Type domain: ")
				return nil
			},
		},
		{
			Name: "Return to previous menu",
		},
	},
	}).Run()
	return domain
}
