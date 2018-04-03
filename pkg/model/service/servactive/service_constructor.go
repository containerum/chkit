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

func RunInteractveConstructor(config ConstructorConfig) (service.Service, error) {
	fmt.Printf("Hi there!\n")
	if !config.Force {
		ok, _ := yes("Do you want to create service?")
		if !ok {
			return service.Service{}, ErrUserStoppedSession
		}
	}
	fmt.Printf("OK\n")
	serv, err := fillServiceField()
	return serv, err
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
		field, ok := askFieldToChange(fields)
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
		case ips:
			IPs, err := getIPs()
			if err != nil {
				return serv, err
			}
			serv.IPs = IPs
		case deploy:
		default:
			panic("[service interactive constructor] unreacheable state in field selection func")
		}
	}
	return serv, nil
}

func defaultService() service.Service {
	return service.Service{
		Name:   namegen.ColoredPhysics(),
		Domain: defaultString,
		IPs:    nil,
		Ports:  nil,
		Deploy: defaultString,
	}
}
