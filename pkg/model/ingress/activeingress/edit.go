package activeingress

import (
	"fmt"

	"strings"

	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/containerum/chkit/pkg/util/tlsview"
	"github.com/sirupsen/logrus"
)

func EditWizard(config Config) (ingress.Ingress, error) {
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
					Label: fmt.Sprintf("Set TLS secret : %s", func() string {
						if rule.TLSSecret == "" {
							return "none"
						}
						return tlsview.SmallView([]byte(rule.TLSSecret))
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
						ingr.Rules = []ingress.Rule{rule}
						data, err := ingr.RenderYAML()
						if err != nil {
							logrus.WithError(err).Errorf("unable to render ingress to yaml")
							activekit.Attention(err.Error())
						}
						border := strings.Repeat("_", text.Width(data))
						fmt.Printf("%s\n%s\n%s\n", border, data, border)
						return nil
					},
				},
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
