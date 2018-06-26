package clingress

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/spf13/cobra"
)

func Delete(ctx *context.Context) *cobra.Command {
	var force = false
	command := &cobra.Command{
		Use:     "ingress",
		Short:   "delete ingress",
		Long:    "Delete ingress.",
		Example: "chkit delete ingress $INGRESS [-n $NAMESPACE] [--force]",
		Aliases: aliases,
		Run: func(cmd *cobra.Command, args []string) {
			logger := coblog.Logger(cmd)
			logger.Debugf("starting ingress delete")
			var ingrName string
			switch len(args) {
			case 0:
				ingrList, err := ctx.Client.GetIngressList(ctx.Namespace.ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get ingress list")
					activekit.Attention("Unable to get ingress list:\n%v", err)
					ctx.Exit(1)
				}
				var menu activekit.MenuItems
				for _, ingr := range ingrList {
					menu = menu.Append(&activekit.MenuItem{
						Label: ingr.Name,
						Action: func(name string) func() error {
							return func() error {
								ingrName = name
								return nil
							}
						}(ingr.Host()),
					})
				}
				(&activekit.Menu{
					Title: "Select ingress",
					Items: menu,
				}).Run()
			case 1:
				name := args[0]
				ingr, err := ctx.Client.GetIngress(ctx.Namespace.ID, name)
				if err != nil {
					logger.WithError(err).Errorf("unable to find ingress %q", name)
					activekit.Attention("Unable to find ingress %q", name)
					ctx.Exit(1)
				}
				ingrName = ingr.Name
			default:
				cmd.Help()
				ctx.Exit(1)
			}
			if force || activekit.YesNo("Do you really want to delete ingress %q?", ingrName) {
				if err := ctx.Client.DeleteIngress(ctx.Namespace.ID, ingrName); err != nil {
					logger.WithError(err).Errorf("unable to delete ingress")
					activekit.Attention("Unable to delete ingress:\n%v", err)
					ctx.Exit(1)
				}
				fmt.Println("Ingress deleted!")
			} else {
				fmt.Println("No ingresses have been deleted")
			}
		},
	}

	command.PersistentFlags().
		BoolVarP(&force, "force", "f", false, "delete ingress without confirmation")

	return command
}
