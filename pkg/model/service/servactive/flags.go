package servactive

import (
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/namegen"
)

type Flags struct {
	Force      bool   `flag:"force f" desc:"suppress confirmation, optional"`
	Name       string `desc:"service name, optional"`
	Deploy     string `desc:"service deployment, required"`
	TargetPort int    `desc:"service target port, optional"`
	Port       int    `desc:"service port, optional"`
	Protocol   string `desc:"service protocol, optional"`
	PortName   string `desc:"service port name, optional"`
}

func (flags Flags) Service() (service.Service, error) {
	var flagSvc = service.Service{
		Name:   flags.Name,
		Deploy: flags.Deploy,
	}

	var flagPort = service.Port{
		Protocol:   flags.Protocol,
		Port:       &flags.Port,
		TargetPort: flags.TargetPort,
	}

	if flags.Name == "" {
		flagSvc.Name = namegen.ColoredPhysics()
	}

	if flags.PortName == "" {
		flagPort.Name = namegen.ColoredPhysics()
	}

	flagSvc.Ports = []service.Port{flagPort}

	return flagSvc, nil
}
