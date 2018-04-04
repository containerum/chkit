package servactive

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/activeToolkit"
	"github.com/containerum/chkit/pkg/util/namegen"
)

const (
	ErrUserExit             chkitErrors.Err = "user exit"
	ErrUserStoppedSession   chkitErrors.Err = "user stopped session"
	ErrInvalidSymbolInLabel chkitErrors.Err = "invalid symbol in label"
	defaultString                           = "undefined"
)

type ConstructorConfig struct {
	Force       bool
	Deployments []string
}

func RunInteractveConstructor(config ConstructorConfig) (service.ServiceList, error) {
	fmt.Printf("Hi there!\n")
	if !config.Force {
		ok, _ := activeToolkit.Yes("Do you want to create service?")
		if !ok {
			return nil, ErrUserExit
		}
		fmt.Printf("OK")
	}
	var list service.ServiceList
	serv := defaultService()
	var err error
	for {
		serv, err = fillServiceField(config, serv)
		switch err {
		case nil:
			// pass
		case ErrUserStoppedSession:
			return list, err
		case ErrUserExit:
			return list, nil
		default:
			return list, err
		}
		if err = validateService(serv); err != nil {
			fmt.Printf("Error: %v", err)
			_, res, _ := activeToolkit.Options("What's next?",
				true,
				"fix service",
				"create new service",
				"exit")
			switch {
			case res == 0:
				continue
			case res == 1:
				serv = defaultService()
				continue
			default:
				return list, ErrUserExit
			}
		}
		if yes, _ := activeToolkit.Yes(fmt.Sprintf("Add %q to list?", serv.Name)); yes {
			list = append(list, serv)
			fmt.Printf("Service %q added to list\n", serv.Name)
		}
		ok, _ := activeToolkit.Yes("Do you want to create another service?")
		if !ok {
			return list, ErrUserStoppedSession
		}
		fmt.Printf("OK")
		serv = defaultService()
	}
}

func fillServiceField(config ConstructorConfig, serv service.Service) (service.Service, error) {
	const (
		name = iota
		deploy
		domain
		ips
		ports
		pushServ
		exit
	)
	for {
		fields := []string{
			fmt.Sprintf("Set name  : %s", serv.Name),
			fmt.Sprintf("Set deploy: %s", serv.Deploy),
			fmt.Sprintf("Set domain: %s", serv.Domain),
			fmt.Sprintf("Set IPs   : [%s]", strings.Join(serv.IPs, ", ")),
			fmt.Sprintf("Set ports : %v", service.PortList(serv.Ports)),
			"Push to list",
			"Exit",
		}
		_, field, _ := activeToolkit.Options("What's next?", false, fields...)
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
			deploy, err := getDeploy(config.Deployments)
			if err != nil {
				return serv, err
			}
			serv.Deploy = deploy
		case pushServ:
			return serv, nil
		case exit:
			return serv, ErrUserExit
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
