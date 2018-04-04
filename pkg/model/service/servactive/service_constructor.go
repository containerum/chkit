package servactive

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/activeToolkit"
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
		ok, _ := activeToolkit.Yes("Do you want to create service?")
		if !ok {
			return nil, ErrUserStoppedSession
		}
		fmt.Printf("OK")
	}
	var list service.ServiceList
	serv := defaultService()
	var err error
	for {
		serv, err = fillServiceField(serv)
		switch err {
		case nil:
			list = append(list, serv)
		case ErrUserStoppedSession:
			ok, _ := activeToolkit.Yes("Do you want to exit")
			if !ok {
				return list, ErrUserStoppedSession
			}
			fmt.Printf("OK")
		default:
			return list, err
		}
		if err = validateService(serv); err != nil {
			fmt.Println(err)
			ok, _ := activeToolkit.Yes("Do you want to fix service?")
			if ok {
				continue
			}
		}
		ok, _ := activeToolkit.Yes("Do you want to create service?")
		if !ok {
			return list, ErrUserStoppedSession
		}
		fmt.Printf("OK")
		serv = defaultService()
	}
	return list, nil
}

func fillServiceField(serv service.Service) (service.Service, error) {
	const (
		name = iota
		deploy
		domain
		ips
		ports
	)
	for {
		fields := []string{
			fmt.Sprintf("Name  : %s", serv.Name),
			fmt.Sprintf("Deploy: %s", serv.Deploy),
			fmt.Sprintf("Domain: %s", serv.Domain),
			fmt.Sprintf("IPs   : [%s]", strings.Join(serv.IPs, ", ")),
			fmt.Sprintf("Ports : %v", service.PortList(serv.Ports)),
		}
		field, ok := activeToolkit.AskFieldToChange(fields)
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
