package cliconfigmap

import (
	"os"

	"fmt"

	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/spf13/cobra"
)

func Get(ctx *context.Context) *cobra.Command {
	var command = &cobra.Command{
		Use:     "configmap",
		Aliases: aliases,
		Run: func(cmd *cobra.Command, args []string) {
			var logger = coblog.Logger(cmd)
			var data model.Renderer
			switch len(args) {
			case 0:
				cm, err := ctx.Client.GetConfigmapList(ctx.Namespace)
				if err != nil {
					logger.WithError(err).Errorf("unable to get configmap list")
					fmt.Printf("Unable to get configmap list:\n%v\n", err)
					os.Exit(1)
				}
				data = cm
			case 1:
				cm, err := ctx.Client.GetConfigmap(ctx.Namespace, args[0])
				if err != nil {
					logger.WithError(err).Errorf("unable to get configmap %q", args[0])
					fmt.Printf("Unable to get configmap %q:\n%v\n", args[0], err)
					os.Exit(1)
				}
				data = cm
			default:
				cmd.Help()
				os.Exit(1)
			}
			var file, _ = cmd.Flags().GetString("file")
			var format, _ = cmd.Flags().GetString("output")
			if err := configuration.ExportData(data, configuration.ExportConfig{
				Filename: file,
				Format:   configuration.ExportFormat(format),
			}); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	var flags = command.PersistentFlags()
	flags.String("file", "-", "output file")
	flags.StringP("output", "o", "", "output format yaml/json")
	return command
}
