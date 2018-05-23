package service

import (
	"fmt"
	"strings"

	"github.com/containerum/kube-client/pkg/model"
)

type Port struct {
	Name       string
	Port       *int
	TargetPort int
	Protocol   string
}

func (port Port) String() string {
	p := fmt.Sprintf("%d", port.TargetPort)
	if port.Port != nil {
		p = fmt.Sprintf("%d:%d", port.TargetPort, *port.Port)
	}
	return fmt.Sprintf("%s %s/%s", port.Name, p, port.Protocol)
}

func PortFromKube(kubePort model.ServicePort) Port {
	return Port{
		Name:       kubePort.Name,
		Port:       kubePort.Port,
		TargetPort: kubePort.TargetPort,
		Protocol:   string(kubePort.Protocol),
	}
}

type PortList []Port

func (list PortList) String() string {
	ports := make([]string, 0, len(list))
	for _, port := range list {
		ports = append(ports, port.String())
	}
	return "[" + strings.Join(ports, ", ") + "]"
}
