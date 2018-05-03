package ingress

import (
	"strings"

	"github.com/containerum/chkit/pkg/model"
)

var (
	_ model.TableRenderer = Ingress{}
	_ model.TableItem     = Ingress{}
)

func (ingress Ingress) RenderTable() string {
	return model.RenderTable(ingress)
}

func (ingress Ingress) TableHeaders() []string {
	return []string{
		"Name",
		"Host",
		"Service",
	}
}

func (ingress Ingress) TableRows() [][]string {
	return [][]string{{
		ingress.Name,
		strings.Join(ingress.Rules.Hosts(), "\n"),
		strings.Join(ingress.Rules.ServicesTableView(), "\n"),
	}}
}

func (ingress Ingress) String() string {
	return ingress.RenderTable()
}
