package servactive

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/validation"
)

func setPortName(port *service.Port) *activekit.MenuItem {
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Set name : %s",
			activekit.OrString(port.Name, "undefined (required)")),
		Action: func() error {
			var promt string
			if port.Name == "" {
				promt = "Print port name: "
			} else {
				promt = fmt.Sprintf("Print port name (hit enter to use %q): ", port.Name)
			}
			name := strings.TrimSpace(activekit.Promt(promt))
			if name == "" {
				return nil
			}
			if err := validation.ValidateLabel(name); err != nil {
				activekit.Attention(fmt.Sprintf("Invalid port name:\n%v", err))
				return nil
			}
			port.Name = name
			return nil
		},
	}
}

func setTargetPort(port *service.Port) *activekit.MenuItem {
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Set target port : %d (required)", port.TargetPort),
		Action: func() error {
			promt := fmt.Sprintf("Print target port (1..65535, hit enter to use %d): ", port.TargetPort)
			portStr := strings.TrimSpace(activekit.Promt(promt))
			if portStr == "" {
				return nil
			}
			var innerPort int
			if _, err := fmt.Sscanf(portStr, "%d", &innerPort); err != nil || (innerPort < 1 && innerPort > 65535) {
				fmt.Printf("Invalid target port %q: must be number 1..65535\n", portStr)
				return nil
			}
			port.TargetPort = innerPort
			return nil
		},
	}
}

func setPortProto(port *service.Port) *activekit.MenuItem {
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Set proto : %s",
			activekit.OrString(port.Protocol, "undefined (required)")),
		Action: func() error {
			_, err := (&activekit.Menu{
				Title: fmt.Sprintf("Select protocol (current: %s)", port.Protocol),
				Items: []*activekit.MenuItem{
					{
						Label: "TCP",
						Action: func() error {
							port.Protocol = "TCP"
							return nil
						},
					},
					{
						Label: "UDP",
						Action: func() error {
							port.Protocol = "UDP"
							return nil
						},
					},
					{
						Label: "Return to previous menu",
					},
				}}).Run()
			return err
		},
	}
}

func setPortPort(port *service.Port) *activekit.MenuItem {
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Set port : %s",
			activekit.OrValue(port.Port, "undefined (optional)")),
		Action: func() error {
			var promt string
			if port.Port == nil {
				promt = "Print port (11000..65535, hit enter to leave empty):"
			} else {
				promt = fmt.Sprintf("Print port (11000..65535, hit enter to use %d, type 'none' to leave empty): ", *port.Port)
			}
			portStr := strings.TrimSpace(activekit.Promt(promt))
			if portStr == "none" {
				port.Port = nil
			}
			if portStr == "" {
				return nil
			}
			var enternalPort int
			if _, err := fmt.Sscanf(portStr, "%d", &enternalPort); err != nil || (enternalPort < 11000 && enternalPort > 65535) {
				fmt.Printf("Invalid port %q: must be number in 11000..65535\n", portStr)
				return nil
			}
			port.Port = &enternalPort
			return nil
		},
	}
}
