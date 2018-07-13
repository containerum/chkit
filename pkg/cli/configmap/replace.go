package cliconfigmap

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/export"
	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/model/configmap/activeconfigmap"
	"github.com/containerum/chkit/pkg/porta"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Replace(ctx *context.Context) *cobra.Command {
	var flags struct {
		activeconfigmap.Flags
		porta.Importer
		porta.Exporter
	}

	exportConfig := export.ExportConfig{}
	var cmd = &cobra.Command{
		Use:     "configmap",
		Short:   "Replace configmap.",
		Aliases: aliases,
		Run: func(cmd *cobra.Command, args []string) {
			var logger = coblog.Logger(cmd)
			logger.Struct(flags)
			logger.Debugf("running replace configmap command")
			var flagCm configmap.ConfigMap

			if flags.ImportActivated() {
				if err := flags.Import(&flagCm); err != nil {
					ferr.Printf("unable to import configmap:\n%v\n", err)
					ctx.Exit(1)
				}
			} else {
				var err error
				flagCm, err = flags.ConfigMap()
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
			}

			cmName := ""
			if flags.Name != "" {
				cmName = flags.Name
			} else if len(args) != 0 {
				cmName = args[0]
			}

			var newCm configmap.ConfigMap
			if flags.Force {
				if cmName == "" {
					cmd.Help()
					return
				}
				oldCm, err := ctx.Client.GetConfigmap(ctx.GetNamespace().ID, cmName)
				if err != nil {
					activekit.Attention(err.Error())
					ctx.Exit(1)
				}
				for k, v := range flagCm.Data {
					oldCm.Data[k] = v
				}

				if err := activeconfigmap.ValidateConfigMap(oldCm); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				if flags.ExporterActivated() {
					if err := flags.Export(oldCm); err != nil {
						ferr.Printf("unable to export configmap:\n%v\n", err)
						ctx.Exit(1)
					}
					return
				}
				if err := ctx.Client.ReplaceConfigmap(ctx.GetNamespace().ID, oldCm.ToBase64()); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Congratulations! Configmap %s updated!\n", oldCm.Name)
				return
			} else {
				if cmName == "" {
					list, err := ctx.Client.GetConfigmapList(ctx.GetNamespace().ID)
					if err != nil {
						activekit.Attention(err.Error())
						ctx.Exit(1)
					}
					var menu []*activekit.MenuItem
					for _, s := range list {
						menu = append(menu, &activekit.MenuItem{
							Label: s.Name,
							Action: func(d configmap.ConfigMap) func() error {
								return func() error {
									newCm = d
									for k, v := range newCm.Data {
										newCm.Data[k] = v
									}
									return nil
								}
							}(s),
						})
					}
					if err := export.ExportData(list, exportConfig); err != nil {
						logrus.WithError(err).Errorf("unable to export data")
						angel.Angel(ctx, err)
					}
					(&activekit.Menu{
						Title: "Choose configmap to replace",
						Items: menu,
					}).Run()
				} else {
					var err error
					newCm, err = ctx.Client.GetConfigmap(ctx.GetNamespace().ID, cmName)
					if err != nil {
						activekit.Attention(err.Error())
						ctx.Exit(1)
					}
				}
			}
			for k, v := range flagCm.Data {
				newCm.Data[k] = v
			}
			if !flags.Force {
				newCm = activeconfigmap.Config{
					EditName:  false,
					ConfigMap: &newCm,
				}.Wizard()
			}
			if flags.Force ||
				activekit.YesNo("Do you really want to replace configmap %q on server?", newCm.Name) {
				if err := ctx.Client.ReplaceConfigmap(ctx.GetNamespace().ID, newCm.ToBase64()); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Congratulations! Configmap %s updated!\n", newCm.Name)
			} else {
				ctx.Exit(0)
			}
			fmt.Println(newCm.RenderTable())
			(&activekit.Menu{
				Items: activekit.MenuItems{
					{
						Label: "Edit configmap " + newCm.Name,
						Action: func() error {
							newCm = activeconfigmap.Config{
								EditName:  false,
								ConfigMap: &newCm,
							}.Wizard()
							if activekit.YesNo("Push changes to server?") {
								if err := ctx.Client.ReplaceConfigmap(ctx.GetNamespace().ID, newCm.ToBase64()); err != nil {
									ferr.Printf("unable to update configmap on server:\n%v\n", err)
								}
							}
							return nil
						},
					},
					{
						Label: "Exit",
						Action: func() error {
							ctx.Exit(0)
							return nil
						},
					},
				},
			}).Run()
		},
	}
	if err := gpflag.ParseTo(&flags, cmd.Flags()); err != nil {
		panic(err)
	}
	return cmd
}
