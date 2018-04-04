package servactive

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/namegen"

	"github.com/containerum/chkit/pkg/chkitErrors"
)

const (
	ErrUserStoppedSession   chkitErrors.Err = "user stopped session"
	ErrInvalidSymbolInLabel chkitErrors.Err = "invalid symbol in label"
	defaultString                           = "undefined"
)

type ConstructorConfig struct {
	Force bool
}

func RunInteractveConstructor(config ConstructorConfig) (service.ServiceList, error) {
	fmt.Printf("Hi there!\n")
	if !config.Force {
		ok, _ := Yes("Do you want to create service?")
		if !ok {
			return nil, ErrUserStoppedSession
		}
		fmt.Printf("OK\n")
	}
	var list service.ServiceList
	serv, err := fillServiceField()
	switch err {
	case ErrUserStoppedSession:
		// pass
	default:
		return nil, err
	}
	list = append(list, serv)
	fmt.Printf("Service %q added to list\n", serv.Name)
	for {
		ok, _ := Yes("Dou you want to create another service?")
		if !ok {
			break
		}
		serv, err = fillServiceField()
		switch err {
		case nil:
			list = append(list, serv)
		case ErrUserStoppedSession:
			continue
		default:
			return list, err
		}
	}
	return list, nil
}

func fillServiceField() (service.Service, error) {
	const (
		name = iota
		domain
		ips
		ports
		deploy
	)
	serv := defaultService()
	for {
		fields := []string{
			fmt.Sprintf("Name  : %s", serv.Name),
			fmt.Sprintf("Domain: %s", serv.Domain),
			fmt.Sprintf("IPs   : [%s]", strings.Join(serv.IPs, ", ")),
			fmt.Sprintf("Ports : %v", service.PortList(serv.Ports)),
			fmt.Sprintf("Deploy: %s", serv.Deploy),
		}
		field, ok := AskFieldToChange(fields)
		if !ok {
			return serv, ErrUserStoppedSession
		}
		switch field {
		case name:
			name, err := getName(serv.Name)
			if err != nil {
				return serv, err
			}
			serv.Name = name
		case ports:
			ports, err := getPorts()
			if err != nil {
				return serv, err
			}
			serv.Ports = ports
		case domain:
			domain, err := getDomain()
			if err != nil {
				return serv, err
			}
			serv.Domain = domain
		case ips:
			IPs, err := getIPs()
			if err != nil {
				return serv, err
			}
			serv.IPs = IPs
		case deploy:
			deploy, err := getDeploy()
			if err != nil {
				return serv, err
			}
			serv.Deploy = deploy
		default:
			panic("[service interactive constructor] unreacheable state in field selection func")
		}
	}
}

func defaultService() service.Service {
	return service.Service{
		Name:   namegen.ColoredPhysics(),
		Domain: "undefined (optional)",
		IPs:    nil,
		Ports:  nil,
		Deploy: "undefined (required)",
	}
}
