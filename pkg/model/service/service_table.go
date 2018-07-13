package service

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"fmt"

	"github.com/containerum/chkit/pkg/model"
	"github.com/ninedraft/boxofstuff/str"
)

var (
	_ model.TableRenderer = new(Service)
)

func (serv Service) RenderTable() string {
	return model.RenderTable(&serv)
}

func (_ *Service) TableHeaders() []string {
	return []string{"Name", "Deploy", "Kind", "Ports", "Age"}
}

func (serv *Service) TableRows() [][]string {
	age := "undefined"
	if serv.CreatedAt != (time.Time{}) {
		age = model.Age(serv.CreatedAt)
	}
	kind := "Internal"

	var ports = make(str.Vector, 0, len(serv.Ports))
	if serv.Domain != "" {
		kind = "External"
		for _, p := range serv.Ports {
			switch strings.ToLower(p.Protocol) {
			case "tcp":
				ports = append(ports, fmt.Sprintf("%s -> %d (%s)", (&url.URL{
					Scheme: "http",
					Host:   serv.Domain + ":" + strconv.Itoa(*p.Port),
				}).String(), p.TargetPort, p.Protocol))
			case "udp":
				ports = append(ports, fmt.Sprintf("%s -> %d (%s)", (&url.URL{
					Scheme: "udp",
					Host:   serv.Domain + ":" + strconv.Itoa(*p.Port),
				}).String(), p.TargetPort, p.Protocol))
			default:
				ports = append(ports, fmt.Sprintf("%s -> %d (%s)", (&url.URL{
					Host: serv.Domain + ":" + strconv.Itoa(*p.Port),
				}).String(), p.TargetPort, p.Protocol))
			}
		}
	} else {
		for _, p := range serv.Ports {
			ports = append(ports, fmt.Sprintf("%5d -> %d (%s)", *p.Port, p.TargetPort, p.Protocol))
		}
	}

	return [][]string{{
		serv.Name,
		str.Vector{serv.Deploy, "!MISSING DEPLOYMENT!"}.FirstNonEmpty(),
		kind,
		strings.Join(ports, "\n"),
		age,
	}}
}
