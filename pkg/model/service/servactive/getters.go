package servactive

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/namegen"
)

func getName(defaultName string) (string, error) {
	for {
		name, _ := AskLine(fmt.Sprintf("Type service name (just leave empty to dub it %s)",
			defaultName))
		if IsStop(name) {
			fmt.Printf("OK :(\n")
			return "", ErrUserStoppedSession
		}
		if name == "" {
			return defaultName, nil
		}
		if err := validateLabel(name); err != nil {
			fmt.Printf("\nError: %v\nPrint new one: ", err)
			continue
		}
		return name, nil
	}

}

func getIPs() ([]string, error) {
	fmt.Printf("Print IP addresses, delimited by spaces or enters.\nPress Ctrl+D or print stop word to end list:\n")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	IPs := make([]string, 0, 4)
	for scanner.Scan() {
		text := scanner.Text()
		if IsStop(text) {
			break
		}
		if net.ParseIP(text) == nil {
			fmt.Printf("\nSorry, but %q is not valid IP address\nPrint new one: ", text)
			continue
		}
		IPs = append(IPs, text)
	}
	return IPs, nil
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
		ok, _ := Yes("Continue creating ports?")
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
		name, _ := AskLine(fmt.Sprintf("type name (hit Enter to use %q) > ", defaultName))
		if IsStop(name) {
			return name, ErrUserStoppedSession
		}
		if name == "" {
			name = defaultName
		}
		if err := validateLabel(name); err != nil {
			fmt.Printf("%v. Try again:\n", err)
			continue
		}
		return name, nil
	}
}

func getPortProtocol(name string) (string, error) {
	for {
		proto, _ := AskLine(fmt.Sprintf("%s::protocol (TCP or UDP , TCP default) > ", name))
		if IsStop(proto) {
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
		target_port_str, exit := AskLine(fmt.Sprintf("%s::target_port > ", name))
		if exit || IsStop(target_port_str) {
			return -1, ErrUserStoppedSession
		}
		target_port, err := strconv.Atoi(target_port_str)
		if err != nil {
			fmt.Printf("Target port can be only number! Try again:\n")
			continue
		}
		return target_port, nil
	}
}

func getOptionalPort(name string) (*int, error) {
	for {
		optional_port_str, exit := AskLine(fmt.Sprintf("%s::port (hit Enter to leave undefined) > ", name))
		if exit || IsStop(optional_port_str) {
			return nil, ErrUserStoppedSession
		}
		if optional_port_str == "" {
			return nil, nil
		}
		optional_port, err := strconv.Atoi(optional_port_str)
		if err != nil {
			fmt.Printf("Port can be only number! Try again:\n")
			continue
		}
		return &optional_port, nil
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
		domain, _ := AskWord("Print domain (hit Ctrl+D or Enter to skip): ")
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

func getDeploy() (string, error) {
	for {
		domain, exit := AskWord("print deploy (hit Ctrl+D or Enter to skip): ")
		if exit {
			return "", ErrUserStoppedSession
		}
		if err := validateLabel(domain); err != nil {
			fmt.Printf("Invalid domain name! Try again.\n")
			continue
		}
		return domain, nil
	}
}
