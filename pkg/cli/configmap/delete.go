package cliconfigmap

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/export"
	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Delete(ctx *context.Context) *cobra.Command {
	var deleteCMConfig = struct {
		Force bool
	}{}
	exportConfig := export.ExportConfig{}
	var command = &cobra.Command{
		Use:     "configmap",
		Short:   "delete configmap",
		Aliases: aliases,
		Run: func(cmd *cobra.Command, args []string) {
			var selectedCM string
			if len(args) == 0 && !deleteCMConfig.Force {
				list, err := ctx.Client.GetConfigmapList(ctx.GetNamespace().ID)
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				if err := export.ExportData(list, exportConfig); err != nil {
					logrus.WithError(err).Errorf("unable to export data")
					angel.Angel(ctx, err)
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
				if len(args) == 0 {
					cmd.Help()
					ctx.Exit(1)
				}
				selectedCM = args[0]
			}
			if deleteCMConfig.Force ||
				activekit.YesNo("Are you sure you want to delete configmap %q in namespace %q?", selectedCM, ctx.GetNamespace()) {
				if err := ctx.Client.DeleteConfigmap(ctx.GetNamespace().ID, selectedCM); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Println("OK")
			}
		},
	}
	command.PersistentFlags().BoolVarP(&deleteCMConfig.Force, "force", "f", false, "delete pod without confirmation")
	return command
}
