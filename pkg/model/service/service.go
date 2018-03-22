package service

import (
	"time"

	"git.containerum.net/ch/kube-client/pkg/model"
)

type Service struct {
	Name      string
	CreatedAt *time.Time
	Deploy    string
	IPs       []string
	Domain    string
	Ports     []Port
}

func ServiceFromKube(kubeService model.Service) Service {
	ports := make([]Port, 0, len(kubeService.Ports))
	for _, kubePort := range kubeService.Ports {
		ports = append(ports, PortFromKube(kubePort))
	}
	var createdAt *time.Time
	if kubeService.CreatedAt != nil {
		t := time.Unix(*kubeService.CreatedAt, 0)
		createdAt = &t
	}
	return Service{
		Name:      kubeService.Name,
		CreatedAt: createdAt,
		Deploy:    kubeService.Deploy,
		IPs:       kubeService.IPs,
		Domain:    kubeService.Domain,
		Ports:     ports,
	}
}
