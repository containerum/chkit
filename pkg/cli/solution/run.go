package clisolution

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/solution"
	"github.com/containerum/chkit/pkg/model/solution/activesolution"
	"github.com/containerum/chkit/pkg/porta"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

var aliases = []string{"sol", "solutions", "sols", "solu", "so"}

func Run(ctx *context.Context) *cobra.Command {
	var flags struct {
		activesolution.Flags
		porta.Importer
		porta.Exporter
	}

	command := &cobra.Command{
		Use:     "solution",
		Aliases: aliases,
		Short:   "run solution from template",
		Example: "chkit run solution [$TEMPLATE] [--env=KEY1:VALUE1,KEY2:VALUE2] [--file $FILENAME] [--force]",
		Run: func(cmd *cobra.Command, args []string) {
			var logger = coblog.Logger(cmd)
			logger.Struct(flags)
			logger.Debugf("running run solution command")
			var sol solution.Solution
			if flags.ImportActivated() {
				if err := flags.Import(&sol); err != nil {
					ferr.Printf("unable to import configmap:\n%v\n", err)
					ctx.Exit(1)
				}
			} else {
				var err error
				sol, err = flags.Solution(ctx.GetNamespace().ID, args)
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
			}
			sol.Namespace = ctx.GetNamespace().ID
			if flags.Force {
				if err := activesolution.ValidateSolution(sol); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				if flags.ExporterActivated() {
					if err := flags.Export(sol); err != nil {
						ferr.Printf("unable to export configmap:\n%v\n", err)
						ctx.Exit(1)
					}
					return
				}
				if err := ctx.GetClient().RunSolution(sol); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Solution %s is ready to run\n", sol.Name)
				return
			}
			solutions, err := ctx.GetClient().GetSolutionsTemplatesList()
			if err != nil {
				ferr.Println(err)
				ctx.Exit(1)
			}
			config := activesolution.WizardConfig{
				EditName:  true,
				Templates: solutions.Names(),
				Solution:  &sol,
			}
			sol = activesolution.Wizard(ctx, config)
			if activekit.YesNo("Are you sure you want to run solution %s?", sol.Name) {
				for k := range sol.Env {
					if sol.Env[k] == "" {
						delete(sol.Env, k)
					}
				}

				if err := ctx.GetClient().RunSolution(sol); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Congratulations! Solution %s is running!\n", sol.Name)
			}
			fmt.Println(sol.RenderTable())
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
