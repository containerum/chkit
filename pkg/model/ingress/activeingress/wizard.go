package activeingress

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/host2dnslabel"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Services service.ServiceList
	Ingress  *ingress.Ingress
}

func Wizard(config Config) (ingress.Ingress, error) {
	var ingr ingress.Ingress
	if config.Ingress != nil {
		ingr = (*config.Ingress).Copy()
	}
	var rule ingress.Rule
	if len(config.Ingress.Rules) > 0 {
		rule = config.Ingress.Rules[0].Copy()
	}
	for exit := false; !exit; {
		_, err := (&activekit.Menu{
			Items: []*activekit.MenuItem{
				{
					Label: fmt.Sprintf("Set host       : %s",
						activekit.OrString(ingr.Host(), "undefined (required)")),
					Action: func() error {
						host := strings.TrimSpace(activekit.Promt(fmt.Sprintf("type host name (hit enter to leave %s): ",
							activekit.OrString(ingr.Host(), ingr.Host()))))
						if host != "" {
							ingr.Name = host2dnslabel.Host2DNSLabel(host)
							rule.Host = host
							ingr.Rules = []ingress.Rule{rule}
						}
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set TLS secret : %s", func() string {
						if rule.TLSSecret == "" {
							return "none"
						}
						return rule.TLSSecret
					}()),
					Action: func() error {
						rule.TLSSecret = tlsSecretMenu(rule.TLSSecret)
						return nil
					},
				},
				{
					Label: "Edit paths",
					Action: func() error {
						rule.Paths = pathsMenu(config.Services, rule.Paths)
						return nil
					},
				},
				{
					Label: "Print to terminal",
					Action: func() error {
						ingr.Rules = ingress.RuleList{rule}
						data, err := ingr.RenderYAML()
						if err != nil {
							logrus.WithError(err).Errorf("unable to render ingress to yaml")
							activekit.Attention(err.Error())
						}
						border := strings.Repeat("_", text.Width(data))
						fmt.Printf("%s\n%s\n%s\n", border, data, border)
						return nil
					},
				}, /* {
					Label: "Save to file",
					Action: func() error {
						logrus.Debugf("saving soltion to file")
						ingr.Rules = ingress.RuleList{rule}
						data, err := ingr.RenderJSON()
						if err != nil {
							logrus.WithError(err).Errorf("unable to render ingress to json")
							activekit.Attention(err.Error())
							return nil
						}
						fname := activekit.Promt("Print filename: ")
						if err := ioutil.WriteFile(fname, []byte(data), os.ModePerm); err != nil {
							logrus.WithError(err).Errorf("unable to write ingress data to file")
							activekit.Attention(err.Error())
							return nil
						}
						fmt.Println("OK")
						return nil
					},
				},*/
				{
					Label: "Confirm",
					Action: func() error {
						ingr.Rules = []ingress.Rule{rule}
						if err := ValidateIngress(ingr); err != nil {
							activekit.Attention(err.Error())
							return nil
						}
						exit = true
						return nil
					},
				},
			},
		}).Run()
		if err != nil {
			return ingr, err
		}
	}
	return ingr, nil
}
