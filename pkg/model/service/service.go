package service

import (
	"time"

	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/model"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Name      string
	CreatedAt *time.Time
	Deploy    string
	IPs       []string
	Domain    string
	Ports     []Port
	origin    *kubeModels.Service
}

func ServiceFromKube(kubeService kubeModels.Service) Service {
	ports := make([]Port, 0, len(kubeService.Ports))
	for _, kubePort := range kubeService.Ports {
		ports = append(ports, PortFromKube(kubePort))
	}
	var createdAt *time.Time
	if kubeService.CreatedAt != nil {
		t, err := time.Parse(model.TimestampFormat, *kubeService.CreatedAt)
		if err != nil {
			logrus.WithError(err).Debugf("invalid created_at timestamp")
		} else {
			createdAt = &t
		}
	}
	return Service{
		Name:      kubeService.Name,
		CreatedAt: createdAt,
		Deploy:    kubeService.Deploy,
		IPs:       kubeService.IPs,
		Domain:    kubeService.Domain,
		Ports:     ports,
		origin:    &kubeService,
	}
}

func (serv *Service) ToKube() kubeModels.Service {
	if serv.origin != nil {
		return *serv.origin
	}
	kubeServ := kubeModels.Service{
		Name:   serv.Name,
		Deploy: serv.Deploy,
		IPs:    serv.IPs,
		Domain: serv.Domain,
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
	serv.origin = &kubeServ
	return *serv.origin
}
