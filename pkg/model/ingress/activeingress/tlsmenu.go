package activeingress

import (
	"fmt"

	"strings"

	"github.com/containerum/chkit/pkg/util/activekit"
)

func tlsSecretMenu(secret *string) *string {
	var ok bool
	var oldSecret *string
	if secret != nil {
		s := *secret
		oldSecret = &s
	}
	for exit := false; !exit; {
		(&activekit.Menu{
			Title: "Edit TLS secret",
			Items: []*activekit.MenuItem{
				{
					Label: fmt.Sprintf("Edit TLS secret name : %s", activekit.OrString(*secret, "undefined, required")),
					Action: func() error {
						scrt := activekit.Promt("Type TLS secret name (hit Enter to leave %s): ", activekit.OrString(*secret, "empty"))
						scrt = strings.TrimSpace(scrt)
						if scrt == "" {
							return nil
						}
						secret = &scrt
						return nil
					},
				},
				{
					Label: "Delete TLS secret",
					Action: func() error {
						if activekit.YesNo("Are you sure you want to delete TLS secret?") {
							secret = nil
						}
						return nil
					},
				},
				{
					Label: "Confirm",
					Action: func() error {
						exit = true
						ok = true
						return nil
					},
				},
				{
					Label: "Return to previous menu, discard all changes",
					Action: func() error {
						exit = true
						ok = false
						return nil
					},
				},
			},
		}).Run()
	}
	if !ok {
		return oldSecret
	}
	return secret
}
