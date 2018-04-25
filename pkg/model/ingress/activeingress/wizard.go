package activeingress

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/tlsview"
)

type Config struct {
	Ingress *ingress.Ingress
}

func Wizard(config Config) (ingress.Ingress, error) {
	var ingr ingress.Ingress
	if config.Ingress != nil {
		ingr = *config.Ingress
	}
	var rule ingress.Rule
	for exit := false; !exit; {
		_, err := (&activekit.Menu{
			Items: []*activekit.MenuItem{
				{
					Label: fmt.Sprintf("Set host       : %q",
						activekit.OrString(ingr.Name, "undefined (required)")),
					Action: func() error {
						host := strings.TrimSpace(activekit.Promt(fmt.Sprintf("type host name (hit enter to leave %s):",
							activekit.OrString(ingr.Name, "empty"))))
						if host != "" {
							ingr.Name = host
							rule.Host = host
						}
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set TLS secret : %s", func() string {
						if rule.TLSSecret == nil {
							return "none"
						}
						return tlsview.SmallView([]byte(*rule.TLSSecret))
					}()),
					Action: func() error {
						rule.TLSSecret = tlsSecretMenu(rule.TLSSecret)
						return nil
					},
				},
				{
					Label: "Confirm",
					Action: func() error {
						// TODO: validation
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
	return ingress.Ingress{}, nil
}
