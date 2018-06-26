package cliconfigmap

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/export"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Get(ctx *context.Context) *cobra.Command {
	var flags struct {
		File   string `desc: "output file"`
		Output string `desc:"output format yaml/json" flag:"output o"`
	}
	var command = &cobra.Command{
		Use:     "configmap",
		Short:   "show configmap data",
		Aliases: aliases,
		Run: func(cmd *cobra.Command, args []string) {
			var logger = coblog.Logger(cmd)
			var data model.Renderer
			switch len(args) {
			case 0:
				cm, err := ctx.GetClient().GetConfigmapList(ctx.GetNamespace().ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get configmap list")
					fmt.Printf("Unable to get configmap list:\n%v\n", err)
					ctx.Exit(1)
				}
				data = cm
			case 1:
				cm, err := ctx.GetClient().GetConfigmap(ctx.GetNamespace().ID, args[0])
				if err != nil {
					logger.WithError(err).Errorf("unable to get configmap %q", args[0])
					fmt.Printf("Unable to get configmap %q:\n%v\n", args[0], err)
					ctx.Exit(1)
				}
				data = cm
			default:
				cmd.Help()
				ctx.Exit(1)
			}
			if err := export.ExportData(data, export.ExportConfig{
				Filename: flags.File,
				Format:   export.ExportFormat(flags.Output),
			}); err != nil {
				ferr.Println(err)
				ctx.Exit(1)
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
