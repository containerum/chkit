package deplactive

import (
	"fmt"

	"git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/text"
)

func editContainerEnvironmentVars(cont *container.Container) {
	envs := cont.Env
	oldEnvs := make([]model.Env, len(cont.Env))
	copy(oldEnvs, cont.Env)
	var ok bool
	const labelWidth = 32
	for exit := false; !exit; {
		var menu []*activekit.MenuItem
		for i, env := range envs {
			menu = append(menu, &activekit.MenuItem{
				Label: fmt.Sprintf("%s:%q", env.Name, text.Crop(env.Value, labelWidth)),
				Action: func(ind int) func() error {
					return func() error {
						environmentVariableMenu(&envs, ind)
						return nil
					}
				}(i),
			})
		}
		(&activekit.Menu{
			Items: append(menu, []*activekit.MenuItem{
				{
					Label: "Add variable",
					Action: func() error {
						env, ok := newEnvVar()
						if ok {
							envs = append(envs, env)
						}
						return nil
					},
				},
				{
					Label: "Confirm",
					Action: func() error {
						ok = true
						exit = true
						return nil
					},
				},
				{
					Label: "Return to previous menu",
					Action: func() error {
						ok = false
						exit = true
						return nil
					},
				},
			}...),
		}).Run()
	}
	if !ok {
		cont.Env = oldEnvs
	} else {
		cont.Env = envs
	}
}

func environmentVariableMenu(envs *[]model.Env, ind int) {
	env := (*envs)[ind]
	for exit := false; !exit; {
		(&activekit.Menu{
			Items: []*activekit.MenuItem{
				{
					Label: fmt.Sprintf("Set name  :  %q", env.Name),
					Action: func() error {
						name := activekit.Promt("Type variable name (hit Enter to leave unchanged): ")
						if name != "" {
							env.Name = name
						}
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set value  :  %q", env.Name),
					Action: func() error {
						value := activekit.Promt("Type variable value (hit Enter to leave unchanged): ")
						if value != "" {
							env.Value = value
						}
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Delete variable %q", env.Name),
					Action: func() error {
						if activekit.YesNo("Are you sure? [Y/N]: ") {
							*envs = append((*envs)[:ind], (*envs)[ind+1:]...)
							exit = true
						}
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
					Label: "Return to previous menu",
					Action: func() error {
						exit = true
						return nil
					},
				},
			},
		}).Run()
	}
}

func newEnvVar() (model.Env, bool) {
	env := model.Env{}
	var ok bool
	for exit := false; !exit; {
		(&activekit.Menu{
			Items: []*activekit.MenuItem{
				{
					Label: fmt.Sprintf("Set name   :  %q",
						activekit.OrString(env.Name, "undefined (required)")),
					Action: func() error {
						name := activekit.Promt("Type variable name (hit Enter to leave unchanged): ")
						if name != "" {
							env.Name = name
						}
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set value  :  %q",
						activekit.OrString(env.Value, "undefined (required)")),
					Action: func() error {
						value := activekit.Promt("Type variable value (hit Enter to leave unchanged): ")
						if value != "" {
							env.Value = value
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
					Label: "Return to previous menu",
					Action: func() error {
						exit = true
						ok = false
						return nil
					},
				},
			},
		}).Run()
	}
	return env, ok
}
