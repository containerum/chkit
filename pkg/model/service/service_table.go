package service

import (
	"net/url"
	"strconv"
	"strings"
	"time"

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
	return []string{"Name", "Deploy", "URL", "Age"}
}

func (serv *Service) TableRows() [][]string {
	age := "undefine"
	if serv.CreatedAt != (time.Time{}) {
		age = model.Age(serv.CreatedAt)
	}
	var links = make(str.Vector, 0, len(serv.Ports))
	if serv.Domain != "" {
		for _, p := range serv.Ports {
			switch strings.ToLower(p.Protocol) {
			case "tcp":
				links = append(links, (&url.URL{
					Scheme: "http",
					Host:   serv.Domain + ":" + strconv.Itoa(*p.Port),
				}).String())
			case "upd":
				links = append(links, (&url.URL{
					Scheme: "udp",
					Host:   serv.Domain + ":" + strconv.Itoa(*p.Port),
				}).String())
			default:
				links = append(links, (&url.URL{
					Host: serv.Domain + ":" + strconv.Itoa(*p.Port),
				}).String())
			}
		}
	}

	return [][]string{{
		serv.Name,
		serv.Deploy,
		strings.Join(links, "\n"),
		age,
	}}
}
