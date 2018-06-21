package container

import (
	"sort"

	"strings"

	"os"

	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/ninedraft/boxofstuff/str"
)

func componentEnvs(cont *container.Container) *activekit.MenuItem {
	var envs = append([]model.Env{}, cont.Env...)
	return &activekit.MenuItem{
		Label: "Edit enviroment variables",
		Action: func() error {
			for exit := false; !exit; {
				sort.Slice(envs, func(i, j int) bool {
					return envs[i].Name < envs[j].Name
				})
				var menuEnvs activekit.MenuItems
				for _, env := range envs {
					menuEnvs = append(menuEnvs, componentEnv(&env, nil))
				}
				menuEnvs = menuEnvs.Append(&activekit.MenuItem{
					Label: "Add env",
					Action: func() error {
						var env = model.Env{}
						var ok = false
						componentEnv(&env, &ok).Action()
						if ok {
							envs = append(envs, env)
						}
						return nil
					},
				},
					&activekit.MenuItem{
						Label: "Confirm",
						Action: func() error {
							cont.Env = envs
							exit = true
							return nil
						},
					},
					&activekit.MenuItem{
						Label: "Drop all changes, return to previous menu",
						Action: func() error {
							exit = true
							return nil
						},
					})
				(&activekit.Menu{
					Title: "Container -> Envs",
					Items: menuEnvs,
				}).Run()
			}
			return nil
		},
	}
}

func componentEnv(oldeEnv *model.Env, ok *bool) *activekit.MenuItem {
	var label = str.Vector{oldeEnv.Name, oldeEnv.Value, "empty env"}.FirstNonEmpty()
	return &activekit.MenuItem{
		Label: "Edit env " + label,
		Action: func() error {
			var env = *oldeEnv
			for exit := false; !exit; {
				(&activekit.Menu{
					Title: "Container -> Envs -> " + label,
					Items: activekit.MenuItems{
						{
							Label: "Edit name",
							Action: func() error {
								var name = activekit.Promt("Type env name, hit Enter to leave %s",
									str.Vector{env.Name, "empty"}.FirstNonEmpty())
								name = strings.TrimSpace(name)
								switch {
								case strings.HasPrefix(name, "$"):
									name = strings.TrimPrefix(name, "$")
									env.Name = name
									env.Value = os.Getenv(name)
								case len(name) > 0:
									env.Name = name
								}
								return nil
							},
						},
						{
							Label: "Edit value",
							Action: func() error {
								var value = activekit.Promt("Type env name, hit Enter to leave %s",
									str.Vector{env.Name, "empty"}.FirstNonEmpty())
								value = strings.TrimSpace(value)
								switch {
								case strings.HasPrefix(value, "$"):
									value = strings.TrimPrefix(value, "$")
									env.Value = os.Getenv(value)
								case len(value) > 0:
									env.Value = value
								}
								return nil
							},
						},
						{
							Label: "Confirm",
							Action: func() error {
								*oldeEnv = env
								exit = true
								if ok != nil {
									*ok = true
								}
								return nil
							},
						},
						{
							Label: "Drop changes, return to previous menu",
							Action: func() error {
								exit = true
								return nil
							},
						},
					},
				}).Run()
			}
			return nil
		},
	}
}
