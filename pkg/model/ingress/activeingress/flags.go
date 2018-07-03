package activeingress

import (
	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/containerum/chkit/pkg/util/host2dnslabel"
	"github.com/containerum/chkit/pkg/util/namegen"
)

type Flags struct {
	Force     bool   `flag:"force f" desc:"suppress confirmation, optional"`
	Name      string `desc:"solution name, optional"`
	Host      string `desc:"ingress host (example: prettyblog.io), required"`
	Service   string `desc:"ingress endpoint service, required"`
	TLSSecret string `desc:"TLS secret string, optional"`
	Path      string `desc:"path to endpoint (example: /content/pages), optional"`
	Port      int    `desc:"ingress endpoint port (example: 80, 443), optional"`
}

func (flags Flags) Ingress() (ingress.Ingress, error) {
	var flagIngress = ingress.Ingress{
		Name: flags.Name,
	}
	var flagRule = ingress.Rule{
		TLSSecret: flags.TLSSecret,
		Host:      flags.Host,
	}
	var flagPath = ingress.Path{
		Path:        flags.Path,
		ServiceName: flags.Service,
		ServicePort: flags.Port,
	}

	if flagPath.Path == "" {
		flagPath.Path = "/"
	}

	if flags.Name == "" {
		flagIngress.Name = namegen.ColoredPhysics()
	}

	if flags.Path != "" ||
		flags.Service != "" ||
		flags.Port != 0 {
		flagRule.Paths = ingress.PathList{flagPath}
	}
	if flags.Host != "" ||
		flags.TLSSecret != "" ||
		flags.Path != "" ||
		flags.Service != "" ||
		flags.Port != 0 {
		flagIngress.Rules = ingress.RuleList{flagRule}
		flagIngress.Name = host2dnslabel.Host2DNSLabel(flagRule.Host)
	}
	return flagIngress, nil
}
