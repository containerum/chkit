package activeingress

import (
	"io/ioutil"
	"path"

	"bytes"
	"encoding/json"
	"os"

	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/containerum/chkit/pkg/util/host2dnslabel"
	"github.com/containerum/chkit/pkg/util/namegen"
	"gopkg.in/yaml.v2"
)

type Flags struct {
	Force     bool   `flag:"force f" desc:"suppress confirmation, optional"`
	File      string `desc:"file with solution data, .yaml or .json, stdin if '-', optional"`
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

	if flags.File != "" {
		var err error
		flagIngress, err = flags.ingressFromFile()
		if err != nil {
			ferr.Println(err)
			return flagIngress, err
		}
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

func (flags Flags) ingressFromFile() (ingress.Ingress, error) {
	var ingr ingress.Ingress
	data, err := func() ([]byte, error) {
		if flags.File == "-" {
			buf := &bytes.Buffer{}
			_, err := buf.ReadFrom(os.Stdin)
			if err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}
		data, err := ioutil.ReadFile(flags.File)
		if err != nil {
			return data, err
		}
		return data, nil
	}()
	if err != nil {
		return ingr, err
	}
	if path.Ext(flags.File) == "yaml" {
		if err := yaml.Unmarshal(data, &ingr); err != nil {
			return ingr, err
		} else {
			if err := json.Unmarshal(data, &ingr); err != nil {
				return ingr, err
			}

		}
	}
	return ingr, nil
}
