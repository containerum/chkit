package cliconfigmap

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/spf13/cobra"
)

func Delete(ctx *context.Context) *cobra.Command {
	var command = &cobra.Command{
		Use:     "configmap",
		Short:   "delete configmap",
		Aliases: aliases,
		Run: func(cmd *cobra.Command, args []string) {
			var selectedCM string
			if len(args) == 0 {
				list, err := ctx.Client.GetConfigmapList(ctx.GetNamespace().ID)
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				var menu activekit.MenuItems
				for _, cm := range list {
					menu = menu.Append(&activekit.MenuItem{
						Label: cm.Name,
						Action: func(cm configmap.ConfigMap) func() error {
							return func() error {
								selectedCM = cm.Name
								return nil
							}
						}(cm.Copy()),
					})
				}
				(&activekit.Menu{
					Title: "Select configmap",
					Items: menu,
				}).Run()
			} else {
				selectedCM = args[0]
			}
			if force, _ := cmd.Flags().GetBool("force"); force ||
				activekit.YesNo("Are you sure you want to delete configmap %q in namespace %q?", selectedCM, ctx.GetNamespace()) {
				if err := ctx.Client.DeleteConfigmap(ctx.GetNamespace().ID, selectedCM); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Println("OK")
			}
		},
	}
	var flags = command.PersistentFlags()
	flags.BoolP("force", "f", false, "suppress confirmation")
	return command
}
