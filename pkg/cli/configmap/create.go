package cliconfigmap

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/model/configmap/activeconfigmap"
	"github.com/containerum/chkit/pkg/porta"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

var aliases = []string{"cm", "confmap", "conf-map", "comap"}

func Create(ctx *context.Context) *cobra.Command {
	var flags struct {
		activeconfigmap.Flags
		porta.Importer
		porta.Exporter
	}

	command := &cobra.Command{
		Use:     "configmap",
		Aliases: aliases,
		Short:   "create configmap",
		Run: func(cmd *cobra.Command, args []string) {
			var logger = coblog.Logger(cmd)
			logger.Struct(flags)
			logger.Debugf("running create configmap command")
			var config configmap.ConfigMap
			if flags.ImportActivated() {
				if err := flags.Import(&config); err != nil {
					ferr.Printf("unable to import configmap:\n%v\n", err)
					ctx.Exit(1)
				}
			} else {
				var err error
				config, err = flags.ConfigMap()
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
			}
			if flags.Force {
				if err := activeconfigmap.ValidateConfigMap(config); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				if flags.ExporterActivated() {
					if err := flags.Export(config); err != nil {
						ferr.Printf("unable to export configmap:\n%v\n", err)
						ctx.Exit(1)
					}
					return
				}
				if err := ctx.Client.CreateConfigMap(ctx.GetNamespace().ID, config.ToBase64()); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Congratulations! ConfigMap %s created!\n", config.Name)
				return
			}
			config = activeconfigmap.Config{
				EditName:  true,
				ConfigMap: &config,
			}.Wizard()
			if activekit.YesNo("Are you sure you want to create configmap %q?", config.Name) {
				if err := ctx.Client.CreateConfigMap(ctx.GetNamespace().ID, config.ToBase64()); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Congratulations! Configmap %s created!\n", config.Name)
			}

			fmt.Println(config.RenderTable())
			(&activekit.Menu{
				Items: activekit.MenuItems{
					{
						Label: "Edit configmap " + config.Name,
						Action: func() error {
							config = activeconfigmap.Config{
								EditName:  false,
								ConfigMap: &config,
							}.Wizard()
							if activekit.YesNo("Push changes to server?") {
								if err := ctx.Client.ReplaceConfigmap(ctx.GetNamespace().ID, config.ToBase64()); err != nil {
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
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
