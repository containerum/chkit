package service

import (
	"git.containerum.net/ch/kube-client/pkg/model"
)

type Port struct {
	Name       string
	Port       *int
	TargetPort int
	Protocol   string
}

func PortFromKube(kubePort model.ServicePort) Port {
	return Port{
		Name:       kubePort.Name,
		Port:       kubePort.Port,
		TargetPort: kubePort.TargetPort,
		Protocol:   string(kubePort.Protocol),
	}
}
