package activesolution

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/context"

	"github.com/containerum/chkit/pkg/model/solution"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/containerum/kube-client/pkg/model"
)

type WizardConfig struct {
	Solution   *solution.Solution
	Namespaces []string
	Templates  []string
	EditName   bool
}

func Wizard(ctx *context.Context, config WizardConfig) solution.Solution {
	var sol = func() solution.Solution {
		if config.Solution != nil {
			return *config.Solution
		}
		return solution.Solution{
			Name:   namegen.Aster() + "-" + namegen.Color(),
			Branch: "master",
		}
	}()

	userEnv := make(map[string]string, 0)
	for k, v := range sol.Env {
		userEnv[k] = v
	}

	for exit := false; !exit; {
		var envItems activekit.MenuItems
		var ind = 0
		for _, env := range sol.EnvironmentVars() {
			envItems = envItems.Append(&activekit.MenuItem{
				Label: fmt.Sprintf("Edit env      : %s", text.Crop(env.String(), 32)),
				Action: func(i int) func() error {
					return func() error {
						env := envMenu(env.ToKube())
						delete(sol.Env, env.Name)
						delete(userEnv, env.Name)
						if env != nil {
							sol.Env[env.Name] = env.Value
							userEnv[env.Name] = env.Value
						} else {
							envItems.Delete(i)
						}
						return nil
					}
				}(ind),
			})
			ind++
		}
		var menu = activekit.MenuItems{
			func() *activekit.MenuItem {
				if config.EditName {
					return &activekit.MenuItem{
						Label: fmt.Sprintf("Edit name     : %s", activekit.OrString(sol.Name, "undefined, required")),
						Action: func() error {
							name := activekit.Promt("Type solution name (hit Enter to leave %s): ", activekit.OrString(sol.Name, "empty"))
							name = strings.TrimSpace(name)
							if name == "" {
								return nil
							}
							sol.Name = name
							return nil
						},
					}
				}
				return nil
			}(),
			{
				Label: fmt.Sprintf("Edit template : %s", activekit.OrString(sol.Template, "undefined, required")),
				Action: func() error {
					var menu activekit.MenuItems
					for _, templ := range config.Templates {
						menu = menu.Append(&activekit.MenuItem{
							Label: templ,
							Action: func(templ string) func() error {
								return func() error {
									sol.Template = templ
									if env, err := ctx.Client.GetSolutionsTemplatesEnvs(sol.Template, sol.Branch); err == nil {
										for k, v := range env.Env {
											if strings.HasPrefix(v, "{{") {
												env.Env[k] = ""
											}
										}
										sol.Env = env.Env
									}
									for k, v := range userEnv {
										sol.Env[k] = v
									}
									return nil
								}
							}(templ),
						})
					}
					(&activekit.Menu{
						Title: "Select template",
						Items: menu.Append(activekit.MenuItems{
							{
								Label: "Return to previous menu",
								Action: func() error {
									exit = true
									return nil
								},
							},
						}...),
					}).Run()
					return nil
				},
			},
			{
				Label: fmt.Sprintf("Edit branch   : %s", activekit.OrString(sol.Branch, "undefined, required")),
				Action: func() error {
					branch := activekit.Promt("Type branch branch (hit Enter to leave %s): ", activekit.OrString(sol.Branch, "empty"))
					branch = strings.TrimSpace(branch)
					if branch == "" {
						return nil
					}
					sol.Branch = branch
					if env, err := ctx.Client.GetSolutionsTemplatesEnvs(sol.Template, sol.Branch); err == nil {
						for k, v := range env.Env {
							if strings.HasPrefix(v, "{{") {
								env.Env[k] = ""
							}
						}
						sol.Env = env.Env
					}
					for k, v := range userEnv {
						sol.Env[k] = v
					}
					return nil
				},
			},
		}
		menu = menu.Append(envItems...).
			Append(&activekit.MenuItem{
				Label: "Add env",
				Action: func() error {
					env := envMenu(model.Env{})
					if env == nil {
						return nil
					}
					sol.Env[env.Name] = env.Value
					userEnv[env.Name] = env.Value
					menu.Append(&activekit.MenuItem{
						Label: fmt.Sprintf("Edit env      : %s", text.Crop(fmt.Sprintf("%s:%q", env.Name, env.Value), 32)),
						Action: func(i int) func() error {
							return func() error {
								env := envMenu(model.Env{
									Name:  env.Name,
									Value: env.Value,
								})
								delete(sol.Env, env.Name)
								delete(userEnv, env.Name)
								if env != nil {
									sol.Env[env.Name] = env.Value
									userEnv[env.Name] = env.Value
								} else {
									envItems.Delete(i)
								}
								return nil
							}
						}(menu.Len()),
					})
					return nil
				},
			}).Append(&activekit.MenuItem{
			Label: "Confirm",
			Action: func() error {
				exit = true
				return nil
			},
		})
		(&activekit.Menu{
			Items: menu,
		}).Run()
	}
	return sol
}
