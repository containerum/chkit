package activeingress

import (
	"fmt"
	"io/ioutil"
	"strings"

	"os"

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
					Label: "Load from file",
					Action: func() error {
						fname := strings.TrimSpace(activekit.Promt("Type filename (hit Enter to return to previous menu): "))
						if fname == "" {
							return nil
						}
						secretBytes, err := ioutil.ReadFile(fname)
						if err != nil {
							activekit.Attention(err.Error())
							return nil
						}
						s := string(secretBytes)
						secret = &s
						return nil
					},
				},
				{
					Label: "Read from stdin",
					Action: func() error {
						if !activekit.YesNo("Are you sure?") {
							return nil
						}
						fmt.Printf("Paste TLS secret here (hit Ctrl+D to submit):\n")
						secretBytes, err := ioutil.ReadAll(os.Stdin)
						if err != nil {
							activekit.Attention(err.Error())
							return nil
						}
						s := string(secretBytes)
						secret = &s
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
					Label: "Return to previous menu, discar all changes",
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
