package cliconfigmap

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/configmap/activeconfigmap"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

var aliases = []string{"cm", "confmap", "conf-map", "comap"}

func Create(ctx *context.Context) *cobra.Command {
	var flags activeconfigmap.Flags
	command := &cobra.Command{
		Use:     "configmap",
		Aliases: aliases,
		Run: func(cmd *cobra.Command, args []string) {
			var logger = coblog.Logger(cmd)
			logger.Struct(flags)
			var config, err = flags.ConfigMap()
			if err != nil {
				ferr.Println(err)
				ctx.Exit(1)
			}
			if !flags.Force {
				config = activeconfigmap.Config{
					EditName:  true,
					ConfigMap: &config,
				}.Wizard()
				fmt.Println(config.RenderTable())
			}
			if flags.Force || activekit.YesNo("Are you sure you want to create configmap %s?", config.Name) {
				if err := config.Validate(); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				if err := ctx.Client.CreateConfigMap(ctx.GetNamespace().ID, config); err != nil {
					logger.WithError(err).Errorf("unable to create configmap %q", config.Name)
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Println("OK")
			} else if !flags.Force {
				config = activeconfigmap.Config{
					EditName:  false,
					ConfigMap: &config,
				}.Wizard()
				fmt.Println(config.RenderTable())
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
