package servactive

import (
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/namegen"
)

type Flags struct {
	Force      bool   `flag:"force f" desc:"suppress confirmation, optional"`
	Name       string `desc:"service name, optional"`
	Deployment string `desc:"service deployment, required"`
	TargetPort int    `desc:"service target port, optional"`
	Port       int    `desc:"service port, optional"`
	Protocol   string `desc:"service protocol, optional"`
	PortName   string `desc:"service port name, optional"`
}

func (flags Flags) Service() (service.Service, error) {
	var flagSvc = service.Service{
		Deploy: flags.Deployment,
	}

	var flagPort = service.Port{
		Protocol:   flags.Protocol,
		TargetPort: flags.TargetPort,
	}

	if flagPort.TargetPort == 0 {
		flagPort.TargetPort = 80
	}

	if flags.Port != 0 {
		flagPort.Port = &flags.Port
	}

	if flags.Protocol != "" {
		flagPort.Protocol = flags.Protocol
	} else {
		flagPort.Protocol = "TCP"
	}

	if flags.Name != "" {
		flagSvc.Name = flags.Name
	} else {
		flagSvc.Name = namegen.ColoredPhysics()
	}

	if flags.PortName != "" {
		flagPort.Name = flags.PortName
	} else {
		flagPort.Name = namegen.ColoredPhysics()
	}

	if flags.Port != 0 || flags.TargetPort != 0 || flags.Protocol != "" || flags.PortName != "" {
		flagSvc.Ports = []service.Port{flagPort}
	} else {
		flagSvc.Ports = []service.Port{}
	}

	return flagSvc, nil
}
