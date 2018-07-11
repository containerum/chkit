package service

import (
	"time"

	"github.com/containerum/chkit/pkg/model"
	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type Service struct {
	Name       string
	CreatedAt  time.Time
	Deploy     string
	IPs        []string
	Domain     string
	SolutionID string
	Ports      []Port
}

func ServiceFromKube(kubeService kubeModels.Service) Service {
	ports := make([]Port, 0, len(kubeService.Ports))
	for _, kubePort := range kubeService.Ports {
		ports = append(ports, PortFromKube(kubePort))
	}
	var createdAt time.Time
	if t, err := time.Parse(model.TimestampFormat, kubeService.CreatedAt); err == nil {
		createdAt = t
	}
	return Service{
		Name:       kubeService.Name,
		CreatedAt:  createdAt,
		Deploy:     kubeService.Deploy,
		IPs:        kubeService.IPs,
		Domain:     kubeService.Domain,
		SolutionID: kubeService.SolutionID,
		Ports:      ports,
	}
}

func (serv *Service) ToKube() kubeModels.Service {
	kubeServ := kubeModels.Service{
		Name:       serv.Name,
		Deploy:     serv.Deploy,
		IPs:        serv.IPs,
		Domain:     serv.Domain,
		SolutionID: serv.SolutionID,
	}
	ports := make([]kubeModels.ServicePort, 0, len(serv.Ports))
	for _, port := range serv.Ports {
		ports = append(ports, kubeModels.ServicePort(kubeModels.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			TargetPort: port.TargetPort,
			Protocol:   kubeModels.Protocol(port.Protocol),
		}))
	}
	kubeServ.Ports = ports
	return kubeServ
}

func (service Service) Copy() Service {
	cp := service
	cp.Ports = append([]Port{}, service.Ports...)
	cp.IPs = append([]string{}, service.IPs...)
	return cp
}

func (service Service) AllTargetPorts() []int {
	ports := make([]int, 0, len(service.Ports))
	for _, port := range service.Ports {
		ports = append(ports, port.TargetPort)
	}
	return ports
}

func (service Service) AllExternalPorts() []int {
	ports := make([]int, 0, len(service.Ports))
	for _, port := range service.Ports {
		if port.Port != nil {
			ports = append(ports, *port.Port)
		}
	}
	return ports
}
