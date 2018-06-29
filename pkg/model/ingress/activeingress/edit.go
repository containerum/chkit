package activeingress

import (
	"fmt"

	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/tlsview"
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
