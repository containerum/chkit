package solution

import (
	"fmt"
	"os"
	"strings"

	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/kube-client/pkg/model"
)

func envMenu(oldEnv model.Env) *model.Env {
	var env = func(e model.Env) *model.Env {
		return &e
	}(oldEnv)
	for exit := false; !exit; {
		(&activekit.Menu{
			Items: activekit.MenuItems{
				{
					Label: fmt.Sprintf("Edit name  : %s", activekit.OrString(env.Name, "undefined, required")),
					Action: func() error {
						var name = activekit.Promt("Type name (hit Enter to leave %s): ", activekit.OrString(oldEnv.Name, "empty"))
						name = strings.TrimSpace(name)
						if name != "" {
							env.Name = name
						}
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Edit value : %s", activekit.OrString(env.Value, "undefined, required")),
					Action: func() error {
						var value = activekit.Promt("Type value (hit Enter to leave %s, start with $ to load value from host environment): ", activekit.OrString(oldEnv.Value, "empty"))
						value = strings.TrimSpace(value)
						if value == "" {
							return nil
						}
						if strings.HasPrefix(value, "$") {
							name := strings.TrimPrefix(value, "$")
							value = os.Getenv(name)
							if env.Name == "" {
								env.Name = name
							}
						}
						env.Value = value
						return nil
					},
				},
				{
					Label: "Confirm",
					Action: func() error {
						exit = true
						return nil
					},
				},
				{
					Label: "Delete env",
					Action: func() error {
						if activekit.YesNo("Are you sure you want to delete env?") {
							env = nil
							exit = true
						}
						return nil
					},
				},
				{
					Label: "Return to previous menu, drop changes",
					Action: func() error {
						env = &oldEnv
						exit = true
						return nil
					},
				},
			},
		}).Run()
	}
	return env
}
