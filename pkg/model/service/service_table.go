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
	return []string{"Name", "Age", "Deploy", "IPs", "Domain", "Ports"}
}

func (serv *Service) TableRows() [][]string {
	ports := make([]string, 0, len(serv.Ports))
	for _, port := range serv.Ports {
		ports = append(ports,
			fmt.Sprintf("%d %s", port.TargetPort, port.Protocol))
	}
	age := "none"
	if serv.CreatedAt != nil {
		age = model.TimestampFormat(*serv.CreatedAt)
	}
	return [][]string{{
		serv.Name,
		age,
		serv.Deploy,
		strings.Join(serv.IPs, "\n"),
		serv.Domain,
		strings.Join(ports, "\n"),
	}}
}
