package service

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = new(Service)
)

func (serv Service) RenderTable() string {
	return model.RenderTable(&serv)
}

func (_ *Service) TableHeaders() []string {
	return []string{"Name", "Deploy", "IPs", "Domain", "Ports", "Age"}
}

func (serv *Service) TableRows() [][]string {
	ports := make([]string, 0, len(serv.Ports))
	for _, port := range serv.Ports {
		optPort := port.TargetPort
		if port.Port != nil {
			optPort = *port.Port
		}
		ports = append(ports,
			fmt.Sprintf("%d:%d/%s", optPort, port.TargetPort, port.Protocol))
	}
	age := "none"
	if serv.CreatedAt != nil {
		age = model.TimestampFormat(*serv.CreatedAt)
	}
	return [][]string{{
		serv.Name,
		serv.Deploy,
		strings.Join(serv.IPs, "\n"),
		serv.Domain,
		strings.Join(ports, "\n"),
		age,
	}}
}
