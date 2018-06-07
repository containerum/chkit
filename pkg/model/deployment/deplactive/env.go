package deplactive

import (
	"fmt"

	"os"
	"strings"

	"io"

	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/containerum/kube-client/pkg/model"
)

func componentEnvs(cont *container.Container) activekit.MenuItems {
	var items = make(activekit.MenuItems, 0, len(cont.Env))
	for _, env := range cont.Env {
		items = items.Append(componentEditEnv(&env))
	}
	return items.Append(&activekit.MenuItem{
		Label: "Add env",
		Action: func() error {
			var env model.Env
			if componentEditEnv(&env).Action() == nil {
				cont.Env = append(cont.Env, env)
			}
			return nil
		},
	})
}

func componentEditEnv(oldEnv *model.Env) *activekit.MenuItem {
	var env = *oldEnv
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Edit env %s:%q", oldEnv.Name, text.Crop(oldEnv.Value, 64)),
		Action: func() error {
			for exit := false; !exit; {
				_, err := (&activekit.Menu{
					Title: "Deployment -> Container -> Env",
					Items: activekit.MenuItems{
						{
							Label: fmt.Sprintf("Edit name : %s",
								activekit.OrString(env.Name, "undefined, required")),
							Action: func() error {
								var name = activekit.Promt("Type env name (hit Enter to leave %s, if has '$' prefix\n"+
									"then try to get value from host env, escape '\\$' to use strictly): ",
									activekit.OrString(env.Name, "empty"))
								name = strings.TrimSpace(name)
								if name != "" {
									switch {
									case strings.HasPrefix(name, "$"):
										env.Name = strings.TrimPrefix(name, "$")
										env.Value = os.Getenv(env.Name)
									case strings.HasPrefix(name, "\\$"):
										env.Name = strings.TrimPrefix(name, "\\")
									default:
										env.Name = name
									}
								}
								return nil
							},
						},
						{
							Label: fmt.Sprintf("Edit value : %q",
								activekit.OrString(env.Value, "undefined, required")),
							Action: func() error {
								var value = activekit.Promt("Type env value (hit Enter to leave %s): ",
									activekit.OrString(env.Value, "empty"))
								value = strings.TrimSpace(value)
								if value != "" {
									env.Value = value
								}
								return nil
							},
						},
						{
							Label: "Confirm",
							Action: func() error {
								if env.Name == "" {
									fmt.Println("name can't be empty!")
								}
								if env.Value == "" {
									fmt.Println("value can't be empty!")
								}
								if env.Name != "" && env.Value != "" {
									*oldEnv = env
									exit = true
								}
								return nil
							},
						},
						{
							Label: "Drop all changes and return to previous menu",
							Action: func() error {
								return io.EOF
							},
						},
					},
				}).Run()
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}
