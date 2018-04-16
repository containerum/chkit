package servactive

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/validation"

	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/namegen"
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
					menu = append(menu[:i], menu[i+1:]...)
					return nil
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
	}
	if ok {
		return ips
	}
	return oldIPs
}

func getPorts() ([]service.Port, error) {
	var ports []service.Port
	for {
		port, err := getPort()
		switch err {
		case nil:
			ports = append(ports, port)
		case ErrUserStoppedSession:
			fmt.Printf("\nPort wasn't added\n")
		default:
			return nil, err
		}
		fmt.Printf("OK, port %q is added\n", port.Name)
		ok, _ := activekit.Yes("Continue creating ports?")
		if !ok {
			break
		}
	}
	fmt.Printf("Added %d ports\n", len(ports))
	return ports, nil
}

func getPort() (service.Port, error) {
	var port service.Port

	name, err := getPortName()
	if err != nil {
		return port, err
	}
	port.Name = name

	proto, err := getPortProtocol(name)
	if err != nil {
		return port, err
	}
	port.Protocol = proto

	target, err := getTargetPort(name)
	if err != nil {
		return port, err
	}
	port.TargetPort = target

	opt, err := getOptionalPort(name)
	if err != nil {
		return port, err
	}
	port.Port = opt

	return port, nil
}

func getPortName() (string, error) {
	for {
		defaultName := namegen.Aster()
		name, _ := activekit.AskLine(fmt.Sprintf("type name (hit Enter to use %q) > ", defaultName))
		if activekit.IsStop(name) {
			return name, ErrUserStoppedSession
		}
		if name == "" {
			name = defaultName
		}
		if err := validation.ValidateLabel(name); err != nil {
			fmt.Printf("%v. Try again:\n", err)
			continue
		}
		return name, nil
	}
}

func getPortProtocol(name string) (string, error) {
	for {
		proto, _ := activekit.AskLine(fmt.Sprintf("%s::protocol (TCP or UDP , TCP default) > ", name))
		if activekit.IsStop(proto) {
			return proto, ErrUserStoppedSession
		}
		switch strings.ToLower(proto) {
		case "tcp", "udp":
		// pass
		case "":
			proto = "TCP"
		default:
			fmt.Printf("Only TCP and UDP protocols are available! You printed %q. Try again:\n", proto)
			continue
		}
		fmt.Printf("Using %s\n", proto)
		return proto, nil
	}
}

func getTargetPort(name string) (int, error) {
	for {
		targetPortStr, exit := activekit.AskLine(fmt.Sprintf("%s::target_port > ", name))
		if exit || activekit.IsStop(targetPortStr) {
			return -1, ErrUserStoppedSession
		}
		targePort, err := strconv.Atoi(targetPortStr)
		if err != nil || targePort < 1 || targePort > 65535 {
			fmt.Printf("Target port can be only number 1..65535! Try again:\n")
			continue
		}
		return targePort, nil
	}
}

func getOptionalPort(name string) (*int, error) {
	for {
		optionalPortStr, exit := activekit.AskLine(fmt.Sprintf("%s::port (hit Enter to leave undefined) > ", name))
		if exit || activekit.IsStop(optionalPortStr) {
			return nil, ErrUserStoppedSession
		}
		if optionalPortStr == "" {
			return nil, nil
		}
		optionalPort, err := strconv.Atoi(optionalPortStr)
		if err != nil || optionalPort < 11000 || optionalPort > 65535 {
			fmt.Printf("Port can be only number 11000..65535! Try again:\n")
			continue
		}
		return &optionalPort, nil
	}
}

func parsePort(text string) (service.Port, error) {
	tokens := strings.Fields(text)
	var input struct {
		Name       string
		Port       string
		TargetPort string
		Protocol   string
	}
	switch len(tokens) {
	case 3:
		input.Name = tokens[0]
		input.Protocol = tokens[1]
		input.TargetPort = tokens[2]
		input.Port = tokens[3]
	}
	return service.Port{}, nil
}

func getDomain() (string, error) {
	for {
		domain, _ := activekit.AskWord("Print domain (hit Ctrl+D or Enter to skip): ")
		if domain == "" {
			return "", nil
		}
		_, err := url.Parse(domain)
		if err != nil {
			fmt.Printf("Invalid domain %q. Try again.\n", domain)
			continue
		}
		return domain, nil
	}
}

func getDeploy(deployments []string) (string, error) {
	for {
		deployment, _, exit := activekit.Options("Choose deployment (print stop to exit):", true, deployments...)
		if exit {
			return "", ErrUserStoppedSession
		}
		return deployment, nil
	}
}
