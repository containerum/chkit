package clisolution

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/spf13/cobra"
)

func Delete(ctx *context.Context) *cobra.Command {
	var force = false
	var command = &cobra.Command{
		Use:     "solution",
		Short:   "Delete running solution",
		Example: "chkit delete solution [--force]",
		Aliases: aliases,
		Run: func(cmd *cobra.Command, args []string) {
			logger := coblog.Logger(cmd)
			logger.Debugf("starting solution delete")
			var solName string
			switch len(args) {
			case 0:
				solList, err := ctx.GetClient().GetRunningSolutionsList(ctx.GetNamespace().ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get solutions list")
					activekit.Attention("Unable to get solutions list:\n%v", err)
					ctx.Exit(1)
				}
				var menu activekit.MenuItems
				for _, sol := range solList.Solutions {
					menu = menu.Append(&activekit.MenuItem{
						Label: sol.Name,
						Action: func(name string) func() error {
							return func() error {
								solName = name
								return nil
							}
						}(sol.Name),
					})
				}
				(&activekit.Menu{
					Title: "Which solution do you want to delete?",
					Items: append(menu, []*activekit.MenuItem{
						{
							Label: "Exit",
						},
					}...),
				}).Run()
				if solName != "" {
					if force || activekit.YesNo("Do you really want to delete solution %q?", solName) {
						if err := ctx.GetClient().DeleteSolution(ctx.GetNamespace().ID, solName); err != nil {
							logger.WithError(err).Errorf("unable to delete solution")
							activekit.Attention("Unable to delete solution:\n%v", err)
							ctx.Exit(1)
						}
						fmt.Println("Solution deleted!")
					} else {
						fmt.Println("No solutions have been deleted")
					}
				}
			case 1:
				name := args[0]
				sol, err := ctx.GetClient().GetRunningSolution(ctx.GetNamespace().ID, name)
				if err != nil {
					logger.WithError(err).Errorf("unable to find solution %q", name)
					activekit.Attention("Unable to find solution %q", name)
					ctx.Exit(1)
				}
				solName = sol.Name
				if force || activekit.YesNo("Do you really want to delete solution %q?", solName) {
					if err := ctx.GetClient().DeleteSolution(ctx.GetNamespace().ID, solName); err != nil {
						logger.WithError(err).Errorf("unable to delete solution")
						activekit.Attention("Unable to delete solution:\n%v", err)
						ctx.Exit(1)
					}
					fmt.Println("Solution deleted!")
				} else {
					fmt.Println("No solutions have been deleted")
				}
			default:
				cmd.Help()
				ctx.Exit(1)
			}
		},
	}
	command.PersistentFlags().
		BoolVarP(&force, "force", "f", false, "delete solution without confirmation")
	return command
}
